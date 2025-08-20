package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"emoji-maker-backend/config"

	"github.com/gofiber/fiber/v2"
)

// 任务状态枚举
const (
	TaskPending   = "PENDING"
	TaskRunning   = "RUNNING"
	TaskSucceeded = "SUCCEEDED"
	TaskFailed    = "FAILED"
	TaskUnknown   = "UNKNOWN"
)

// 视频生成请求体
type VideoCreateRequest struct {
	ImgBase64      string `json:"img_base64"`      // 图片Base64编码 (图生视频)
	Type           string `json:"type"`            // 生成类型: text_to_video 或 image_to_video
	Prompt         string `json:"prompt"`          // 核心描述文本
	NegativePrompt string `json:"negative_prompt"` // 反向提示词
	Size           string `json:"size"`            // 视频分辨率 (文生视频)
	Resolution     string `json:"resolution"`      // 视频分辨率档位 (图生视频)
}

// 创建任务响应
type CreateTaskResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		JobID string `json:"job_id"`
	} `json:"data"`
}

// 查询任务结果响应
type QueryTaskResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    struct {
		JobID        string `json:"job_id"`
		Status       string `json:"status"`
		VideoURL     string `json:"video_url,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
	} `json:"data"`
}

// DashScope API请求体
type DashScopeRequest struct {
	Model      string `json:"model"`
	Input      Input  `json:"input"`
	Parameters Params `json:"parameters,omitempty"`
}

type Input struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	ImgURL         string `json:"img_url,omitempty"`
}

type Params struct {
	Size         string `json:"size,omitempty"`
	Resolution   string `json:"resolution,omitempty"`
	Duration     int    `json:"duration,omitempty"`
	PromptExtend bool   `json:"prompt_extend,omitempty"`
	Seed         int    `json:"seed,omitempty"`
	Watermark    bool   `json:"watermark,omitempty"`
}

// DashScope API响应体
type DashScopeResponse struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
	} `json:"output"`
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// DashScope 查询任务结果响应
type DashScopeQueryResponse struct {
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID        string `json:"task_id"`
		TaskStatus    string `json:"task_status"`
		SubmitTime    string `json:"submit_time"`
		ScheduledTime string `json:"scheduled_time"`
		EndTime       string `json:"end_time"`
		VideoURL      string `json:"video_url,omitempty"`
		OrigPrompt    string `json:"orig_prompt"`
		ActualPrompt  string `json:"actual_prompt,omitempty"`
	} `json:"output"`
	Usage struct {
		Duration   int    `json:"duration"`
		VideoCount int    `json:"video_count"`
		SR         int    `json:"SR,omitempty"`
		VideoRatio string `json:"video_ratio,omitempty"`
	} `json:"usage"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// 生成随机任务ID
func generateJobID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "job_" + hex.EncodeToString(bytes), nil
}

// 处理用户输入的提示词，通过文本模型优化提示词
func processPromptWithTextModel(originalPrompt string, systemPrompt string) (string, error) {
	apiKey := config.AppConfig.AI.Key
	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "qwen-flash",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": originalPrompt,
			},
		},
		"enable_search": true,
		"forced_search": true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return originalPrompt, fmt.Errorf("failed to marshal request body: %w", err)
	}

	cmd := exec.Command("curl", "-X", "POST", url,
		"-H", "Authorization: Bearer "+apiKey,
		"-H", "Content-Type: application/json",
		"-d", string(jsonData),
		"-s") // 添加 -s 来禁止进度条

	output, err := cmd.CombinedOutput()
	if err != nil {
		return originalPrompt, fmt.Errorf("curl command failed: %w, output: %s", err, string(output))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	fmt.Println("DashScope API Response:", string(output))
	if err := json.Unmarshal(output, &response); err != nil {
		return originalPrompt, fmt.Errorf("failed to unmarshal response: %w, response: %s", err, string(output))
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	return originalPrompt, fmt.Errorf("no content in response")
}

// 创建视频生成任务
func CreateVideoTask(c *fiber.Ctx) error {
	// 从 form-data 中解析字段
	req := VideoCreateRequest{
		Type:           c.FormValue("type"),
		Prompt:         c.FormValue("prompt"),
		NegativePrompt: c.FormValue("negative_prompt"),
		Size:           c.FormValue("size"),
		Resolution:     c.FormValue("resolution"),
		ImgBase64:      c.FormValue("img_base64"),
	}

	// 验证必填字段
	if req.Type != "text_to_video" && req.Type != "image_to_video" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid type. Must be 'text_to_video' or 'image_to_video'",
		})
	}

	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	// 根据类型验证其他字段
	if req.Type == "text_to_video" && req.Size == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Size is required for text_to_video",
		})
	}

	if req.Type == "image_to_video" {
		if req.Resolution == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Resolution is required for image_to_video",
			})
		}
		if req.ImgBase64 == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "img_base64 is required for image_to_video",
			})
		}
	}

	// 生成任务ID
	jobID, err := generateJobID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate job ID",
		})
	}

	// 调用DashScope API创建任务
	var dashScopeReq DashScopeRequest

	// 根据类型设置模型和参数
	if req.Type == "text_to_video" {
		// 文生视频使用 wanx2.1-t2v-turbo 模型
		dashScopeReq = DashScopeRequest{
			Model: "wanx2.1-t2v-turbo",
			Input: Input{
				Prompt:         req.Prompt,
				NegativePrompt: req.NegativePrompt,
			},
			Parameters: Params{
				Size: req.Size,
			},
		}
	} else if req.Type == "image_to_video" {
		// 图生视频使用 wan2.2-i2v-flash 模型
		dashScopeReq = DashScopeRequest{
			Model: "wan2.2-i2v-flash",
			Input: Input{
				Prompt:         req.Prompt,
				NegativePrompt: req.NegativePrompt,
				ImgURL:         req.ImgBase64,
			},
			Parameters: Params{
				Resolution: req.Resolution,
			},
		}
	}

	// 保存任务信息到文件(模拟任务队列/数据库)
	taskFile := fmt.Sprintf("tasks/%s.json", jobID)
	os.MkdirAll("tasks", 0755)

	taskData := map[string]interface{}{
		"job_id":            jobID,
		"status":            TaskPending,
		"request":           req,
		"created_at":        time.Now().Format(time.RFC3339),
		"dashscope_task_id": "",
	}

	taskJSON, _ := json.Marshal(taskData)
	os.WriteFile(taskFile, taskJSON, 0644)

	// 异步调用DashScope API
	go func() {
		// 更新任务状态为运行中
		taskData["status"] = TaskRunning
		taskJSON, _ := json.Marshal(taskData)
		os.WriteFile(taskFile, taskJSON, 0644)

		// 调用DashScope API
		apiKey := config.AppConfig.AI.Key

		url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis"

		jsonData, _ := json.Marshal(dashScopeReq)

		// 创建HTTP请求
		request, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
		request.Header.Set("Authorization", "Bearer "+apiKey)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-DashScope-Async", "enable")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			// 更新任务状态为失败
			taskData["status"] = TaskFailed
			taskData["error"] = err.Error()
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		var dashScopeResp DashScopeResponse
		if err := json.Unmarshal(body, &dashScopeResp); err != nil {
			// 更新任务状态为失败
			taskData["status"] = TaskFailed
			taskData["error"] = "Failed to parse DashScope response: " + err.Error() + ", response: " + string(body)
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}

		if dashScopeResp.Code != "" {
			// 更新任务状态为失败
			taskData["status"] = TaskFailed
			taskData["error"] = dashScopeResp.Message
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}

		// 保存DashScope任务ID
		taskData["dashscope_task_id"] = dashScopeResp.Output.TaskID
		taskData["status"] = TaskRunning
		taskJSON, _ = json.Marshal(taskData)
		os.WriteFile(taskFile, taskJSON, 0644)
	}()

	// 返回成功响应
	response := CreateTaskResponse{
		Code:    200,
		Message: "任务创建成功",
	}
	response.Data.JobID = jobID

	return c.JSON(response)
}

// 查询任务结果
func GetVideoTaskResult(c *fiber.Ctx) error {
	jobID := c.Params("job_id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Job ID is required",
		})
	}

	// 读取任务文件
	taskFile := fmt.Sprintf("tasks/%s.json", jobID)
	taskDataBytes, err := os.ReadFile(taskFile)
	if err != nil {
		// 任务不存在
		response := QueryTaskResponse{
			Code:    500,
			Message: "任务不存在",
		}
		response.Data.JobID = jobID
		response.Data.Status = TaskUnknown
		return c.JSON(response)
	}

	var taskData map[string]interface{}
	if err := json.Unmarshal(taskDataBytes, &taskData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse task data",
		})
	}

	status := taskData["status"].(string)
	dashscopeTaskID, _ := taskData["dashscope_task_id"].(string)

	// 如果是运行中或等待中的任务，并且已经获取到了dashscope_task_id，就查询DashScope API获取最新状态
	if (status == TaskRunning || status == TaskPending) && dashscopeTaskID != "" {
		// 查询DashScope任务状态
		apiKey := config.AppConfig.AI.Key

		url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", dashscopeTaskID)

		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			// 查询失败，将任务标记为失败并记录错误
			status = TaskFailed
			taskData["status"] = status
			taskData["error"] = "Failed to query DashScope task status: " + err.Error()
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
		} else {
			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)

			var dashScopeQueryResp DashScopeQueryResponse
			if err := json.Unmarshal(body, &dashScopeQueryResp); err == nil {
				// 更新本地任务状态
				switch dashScopeQueryResp.Output.TaskStatus {
				case "PENDING":
					status = TaskPending
				case "RUNNING":
					status = TaskRunning
				case "SUCCEEDED":
					status = TaskSucceeded
					taskData["video_url"] = dashScopeQueryResp.Output.VideoURL
				case "FAILED":
					status = TaskFailed
					// 优先使用DashScope返回的详细错误信息
					if dashScopeQueryResp.Message != "" {
						taskData["error"] = dashScopeQueryResp.Message
					} else {
						taskData["error"] = "Task failed on DashScope without a specific message."
					}
				case "UNKNOWN":
					status = TaskUnknown
				default:
					status = TaskUnknown
				}

				// 更新任务文件
				taskData["status"] = status
				taskJSON, _ := json.Marshal(taskData)
				os.WriteFile(taskFile, taskJSON, 0644)
			} else {
				// JSON解析失败，也标记为失败
				status = TaskFailed
				taskData["status"] = status
				taskData["error"] = "Failed to parse DashScope query response: " + err.Error() + ", response body: " + string(body)
				taskJSON, _ := json.Marshal(taskData)
				os.WriteFile(taskFile, taskJSON, 0644)
			}
		}
	}

	// 如果任务成功，将视频转换为GIF
	if status == TaskSucceeded {
		if videoURL, ok := taskData["video_url"].(string); ok && !strings.HasSuffix(videoURL, ".gif") {
			fmt.Println("开始下载视频:", videoURL)
			// 1. 下载视频文件
			videoPath := fmt.Sprintf("tasks/%s.mp4", jobID)
			resp, err := http.Get(videoURL)
			if err != nil {
				status = TaskFailed
				taskData["status"] = status
				taskData["error"] = "Failed to download video: " + err.Error()
				fmt.Println("下载视频失败:", err)
			} else {
				defer resp.Body.Close()
				out, err := os.Create(videoPath)
				if err != nil {
					status = TaskFailed
					taskData["status"] = status
					taskData["error"] = "Failed to create video file: " + err.Error()
					fmt.Println("创建视频文件失败:", err)
				} else {
					defer out.Close()
					_, err = io.Copy(out, resp.Body)
					if err != nil {
						status = TaskFailed
						taskData["status"] = status
						taskData["error"] = "Failed to save video file: " + err.Error()
						fmt.Println("保存视频文件失败:", err)
					} else {
						fmt.Println("视频下载成功:", videoPath)
						// 2. 本地转换为标准GIF (适合微信发送的尺寸)
						gifPath := fmt.Sprintf("tasks/%s.gif", jobID)
						// 使用两步法优化GIF，保持原始宽高比:
						// 第一步: 生成调色板
						palettePath := fmt.Sprintf("tasks/%s_palette.png", jobID)
						paletteCmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "scale=240:240:force_original_aspect_ratio=decrease,pad=240:240:(ow-iw)/2:(oh-ih)/2:color=black@0,palettegen", palettePath)
						paletteOutput, err := paletteCmd.CombinedOutput()
						if err != nil {
							status = TaskFailed
							taskData["status"] = status
							taskData["error"] = "Failed to generate palette: " + err.Error()
							fmt.Println("调色板生成失败:", err, string(paletteOutput))
						} else {
							// 第二步: 使用调色板生成优化的GIF，保持原始宽高比
							cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", palettePath, "-lavfi", "scale=240:240:force_original_aspect_ratio=decrease,pad=240:240:(ow-iw)/2:(oh-ih)/2:color=black@0,fps=8 [x]; [x][1:v] paletteuse", "-f", "gif", gifPath)
							output, err := cmd.CombinedOutput()
							if err != nil {
								status = TaskFailed
								taskData["status"] = status
								taskData["error"] = "Failed to convert video to GIF: " + err.Error()
								fmt.Println("ffmpeg转换失败:", err, string(output))
							} else {
								finalGifURL := fmt.Sprintf("https://"+config.AppConfig.Server.Host+":"+config.AppConfig.Server.Port+"/tasks/%s.gif", jobID)
								taskData["video_url"] = finalGifURL
								fmt.Println("GIF生成成功:", finalGifURL)
							}
							// 清理调色板文件
							os.Remove(palettePath)
						}
						// 3. 清理临时视频文件
						os.Remove(videoPath)
					}
				}
			}
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
		}
	}

	// 构造响应
	response := QueryTaskResponse{
		Code: 200,
	}
	response.Data.JobID = jobID
	response.Data.Status = status

	if status == TaskSucceeded {
		if videoURL, ok := taskData["video_url"].(string); ok {
			response.Data.VideoURL = videoURL
		}
	} else if status == TaskFailed {
		if errorMsg, ok := taskData["error"].(string); ok {
			response.Data.ErrorMessage = errorMsg
		} else {
			response.Data.ErrorMessage = "视频生成失败"
		}
	}

	return c.JSON(response)
}

// VideoCreateRequestWithPromptProcessing defines the request for the new video creation endpoint
type VideoCreateRequestWithPromptProcessing struct {
	Role   string `json:"role"`
	Source string `json:"source"` // Optional
	Action string `json:"action"`
	Size   string `json:"size"`
}

// CreateVideoTaskWithPromptProcessing handles the new video creation process
func CreateVideoTaskWithPromptProcessing(c *fiber.Ctx) error {
	// 1. Parse request from form-data
	req := VideoCreateRequestWithPromptProcessing{
		Role:   c.FormValue("role"),
		Source: c.FormValue("source"),
		Action: c.FormValue("action"),
		Size:   c.FormValue("size"),
	}

	if req.Role == "" || req.Action == "" || req.Size == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "role, action, and size are required"})
	}

	// 2. First prompt processing
	roleInfo := req.Role
	if req.Source != "" {
		roleInfo = fmt.Sprintf("%s (来自 %s)", req.Role, req.Source)
	}

	processedPrompt1, err := processPromptWithTextModel(roleInfo, "你是一个角色描绘大师，请你根据网络搜索使用语言详细的描述这个角色的外貌，体型，颜色，外观,请开始分析并严格按照格式输出，只输出最终的描述，不要包含任何markdown格式或标题。")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed in first prompt processing step: " + err.Error()})
	}

	// 4. Create video task
	jobID, err := generateJobID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate job ID"})
	}

	dashScopeReq := DashScopeRequest{
		Model: "wanx2.1-t2v-turbo",
		Input: Input{
			Prompt: fmt.Sprintf("角色:%s。角色描述: %s。动作: %s。", roleInfo, processedPrompt1, req.Action),
		},
		Parameters: Params{
			Size: req.Size,
		},
	}

	taskFile := fmt.Sprintf("tasks/%s.json", jobID)
	os.MkdirAll("tasks", 0755)

	taskData := map[string]interface{}{
		"job_id":            jobID,
		"status":            TaskPending,
		"request":           req,
		"created_at":        time.Now().Format(time.RFC3339),
		"dashscope_task_id": "",
		"final_prompt":      fmt.Sprintf("角色:%s。角色描述: %s。动作: %s。", roleInfo, processedPrompt1, req.Action),
	}
	taskJSON, _ := json.Marshal(taskData)
	os.WriteFile(taskFile, taskJSON, 0644)

	go func() {
		taskData["status"] = TaskRunning
		taskJSON, _ := json.Marshal(taskData)
		os.WriteFile(taskFile, taskJSON, 0644)

		apiKey := config.AppConfig.AI.Key
		url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis"
		jsonData, _ := json.Marshal(dashScopeReq)

		request, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
		request.Header.Set("Authorization", "Bearer "+apiKey)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-DashScope-Async", "enable")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			taskData["status"] = TaskFailed
			taskData["error"] = err.Error()
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		var dashScopeResp DashScopeResponse
		if err := json.Unmarshal(body, &dashScopeResp); err != nil {
			taskData["status"] = TaskFailed
			taskData["error"] = "Failed to parse DashScope response: " + err.Error() + ", response: " + string(body)
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}

		if dashScopeResp.Code != "" {
			taskData["status"] = TaskFailed
			taskData["error"] = dashScopeResp.Message
			taskJSON, _ := json.Marshal(taskData)
			os.WriteFile(taskFile, taskJSON, 0644)
			return
		}

		taskData["dashscope_task_id"] = dashScopeResp.Output.TaskID
		taskData["status"] = TaskRunning
		taskJSON, _ = json.Marshal(taskData)
		os.WriteFile(taskFile, taskJSON, 0644)
	}()

	response := CreateTaskResponse{
		Code:    200,
		Message: "任务创建成功",
	}
	response.Data.JobID = jobID

	return c.JSON(response)
}
