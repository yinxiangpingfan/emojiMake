package com.example.emojimake.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.emojimake.network.QueryTaskResponseData
import com.example.emojimake.repository.VideoRepository
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import com.example.emojimake.network.BaseResponse
import com.google.gson.Gson
import kotlinx.coroutines.launch
import retrofit2.HttpException

sealed class VideoResult {
    data class Success(val data: QueryTaskResponseData) : VideoResult()
    data class Error(val message: String) : VideoResult()
    object Loading : VideoResult()
    object Idle : VideoResult()
}

class VideoViewModel(private val videoRepository: VideoRepository) : ViewModel() {

    private val _videoResult = MutableStateFlow<VideoResult>(VideoResult.Idle)
    val videoResult: StateFlow<VideoResult> = _videoResult.asStateFlow()
    
    // 添加一个状态来通知token过期
    private val _tokenExpired = MutableStateFlow(false)
    val tokenExpired: StateFlow<Boolean> = _tokenExpired.asStateFlow()

    private fun handleVideoError(response: BaseResponse<*>) {
        val errorMessage = when (response.getErrorMessage()) {
            "Invalid or expired JWT" -> {
                // 当token过期时，设置tokenExpired为true
                _tokenExpired.value = true
                "登录状态无效或已过期"
            }
            else -> response.getErrorMessage()
        }
        _videoResult.value = VideoResult.Error(errorMessage)
    }

    fun createTextToVideo(prompt: String, negativePrompt: String?) {
        viewModelScope.launch {
            _videoResult.value = VideoResult.Loading
            try {
                val response = videoRepository.createVideo(
                    type = "text_to_video",
                    prompt = prompt,
                    negativePrompt = negativePrompt,
                    size = "624*624",
                    resolution = null,
                    imgBase64 = null
                )
                if ((response.code == 200 || response.code == 0) && response.data != null) {
                    startPolling(response.data.job_id)
                } else {
                    handleVideoError(response)
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                handleVideoError(errorResponse)
            } catch (e: Exception) {
                _videoResult.value = VideoResult.Error(e.message ?: "未知错误")
            }
        }
    }

    fun createImageToVideo(prompt: String, negativePrompt: String?, imgBase64: String) {
        viewModelScope.launch {
            _videoResult.value = VideoResult.Loading
            try {
                val response = videoRepository.createVideo(
                    type = "image_to_video",
                    prompt = prompt,
                    negativePrompt = negativePrompt,
                    size = null,
                    resolution = "480P",
                    imgBase64 = imgBase64
                )
                if ((response.code == 200 || response.code == 0) && response.data != null) {
                    startPolling(response.data.job_id)
                } else {
                    handleVideoError(response)
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                handleVideoError(errorResponse)
            } catch (e: Exception) {
                _videoResult.value = VideoResult.Error(e.message ?: "未知错误")
            }
        }
    }

    fun createVideoWithPrompt(role: String, source: String?, action: String) {
        viewModelScope.launch {
            _videoResult.value = VideoResult.Loading
            try {
                val response = videoRepository.createVideoWithPrompt(
                    role = role,
                    source = source,
                    action = action,
                    size = "624*624"
                )
                if ((response.code == 200 || response.code == 0) && response.data != null) {
                    startPolling(response.data.job_id)
                } else {
                    handleVideoError(response)
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                handleVideoError(errorResponse)
            } catch (e: Exception) {
                _videoResult.value = VideoResult.Error(e.message ?: "未知错误")
            }
        }
    }

    private fun startPolling(jobId: String) {
        viewModelScope.launch {
            while (true) {
                try {
                    val response = videoRepository.queryTask(jobId)
                    if ((response.code == 200 || response.code == 0) && response.data != null) {
                        when (response.data.status) {
                            "SUCCEEDED" -> {
                                _videoResult.value = VideoResult.Success(response.data)
                                break
                            }
                            "FAILED" -> {
                                _videoResult.value = VideoResult.Error(response.data.error_message ?: "任务失败")
                                break
                            }
                        }
                    } else {
                        if (response.data?.status == "UNKNOWN" || response.code == 500) {
                            _videoResult.value = VideoResult.Error("任务不存在")
                        } else {
                            handleVideoError(response)
                        }
                        break
                    }
                } catch (e: HttpException) {
                    val errorBody = e.response()?.errorBody()?.string()
                    val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                    handleVideoError(errorResponse)
                    break
                } catch (e: Exception) {
                    _videoResult.value = VideoResult.Error(e.message ?: "未知错误")
                    break
                }
                delay(2500) // Poll every 2.5 seconds
            }
        }
    }

    fun resetVideoResult() {
        _videoResult.value = VideoResult.Idle
    }
    
    // 添加重置token过期状态的方法
    fun resetTokenExpired() {
        _tokenExpired.value = false
    }
}