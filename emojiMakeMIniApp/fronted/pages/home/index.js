// pages/home/index.js
const app = getApp();
const api = require('../../utils/api.js');

Page({
  data: {
    currentMode: 'image_to_video', // 'image_to_video', 'text_to_video' 或 'text_to_video_online'
    tempFilePath: '',
    prompt: '',
    negativePrompt: '',
    loading: false,
    resultVideoUrl: '',
    jobId: '',
    pollingTimer: null,
    pollingCount: 0,
    imageMimeType: 'image/png', // 默认MIME类型
    
    // 新增的文生视频（联网）相关字段
    role: '',
    source: '',
    action: '',
  },

  onLoad: function () {
    // 页面加载时不需要处理用户信息
  },

  chooseImage: function () {
    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFilePaths[0];
        // 尝试获取文件MIME类型
        const fileExtension = tempFilePath.split('.').pop().toLowerCase();
        let mimeType = 'image/jpeg'; // 默认为jpeg
        if (fileExtension === 'png') {
          mimeType = 'image/png';
        } else if (fileExtension === 'gif') {
          mimeType = 'image/gif';
        }
        
        this.setData({
          tempFilePath: tempFilePath,
          imageMimeType: mimeType
        });
      }
    });
  },

  onPromptInput: function (e) {
    this.setData({
      prompt: e.detail.value
    });
  },

  onNegativePromptInput: function (e) {
    this.setData({
      negativePrompt: e.detail.value
    });
  },
  
  // 新增的输入处理函数
  onRoleInput: function (e) {
    this.setData({
      role: e.detail.value
    });
  },
  
  onSourceInput: function (e) {
    this.setData({
      source: e.detail.value
    });
  },
  
  onActionInput: function (e) {
    this.setData({
      action: e.detail.value
    });
  },

  switchMode: function(e) {
    this.setData({
      currentMode: e.currentTarget.dataset.mode
    })
  },

  bindSizeChange: function(e) {
    this.setData({
      sizeIndex: e.detail.value
    })
  },

  bindResolutionChange: function(e) {
    this.setData({
      resolutionIndex: e.detail.value
    })
  },

  generateVideo: function () {
    // 统一校验逻辑
    if (this.data.currentMode === 'text_to_video_online') {
      if (!this.data.role || !this.data.action) {
        wx.showToast({ title: '请输入角色和动作', icon: 'none' });
        return;
      }
    } else if (this.data.currentMode === 'text_to_video') {
      if (!this.data.prompt) {
        wx.showToast({ title: '请输入描述', icon: 'none' });
        return;
      }
    } else if (this.data.currentMode === 'image_to_video') {
      if (!this.data.prompt) {
        wx.showToast({ title: '请输入描述', icon: 'none' });
        return;
      }
      if (!this.data.tempFilePath) {
        wx.showToast({ title: '请先选择图片', icon: 'none' });
        return;
      }
    }

    this.setData({ loading: true, resultVideoUrl: '' });

    let taskType = this.data.currentMode;
    let params = {};
    
    // 根据不同的模式设置参数，严格按照接口文档要求
    if (taskType === 'text_to_video_online') {
      // 文生视频（联网）模式 - 使用 /api/v1/video/create_with_prompt 接口
      // 必需参数：role, action, size
      // 可选参数：source
      params = {
        role: this.data.role,
        action: this.data.action,
        size: "624*624"
      };
      
      // 只有当source有值时才添加
      if (this.data.source && this.data.source.trim()) {
        params.source = this.data.source;
      }
    } else if (taskType === 'text_to_video') {
      // 普通文生视频模式 - 使用 /api/v1/video/create 接口
      // 必需参数：type, prompt, size (当type为text_to_video时)
      // 可选参数：negative_prompt
      params = {
        type: 'text_to_video',
        prompt: this.data.prompt,
        size: "624*624"
      };
      
      // 只有当negative_prompt有值时才添加
      if (this.data.negativePrompt && this.data.negativePrompt.trim()) {
        params.negative_prompt = this.data.negativePrompt;
      }
    } else if (taskType === 'image_to_video') {
      // 图生视频模式 - 使用 /api/v1/video/create 接口
      // 必需参数：type, prompt, resolution (当type为image_to_video时), img_base64
      // 可选参数：negative_prompt
      params = {
        type: 'image_to_video',
        prompt: this.data.prompt,
        resolution: "480P"
      };
      
      // 只有当negative_prompt有值时才添加
      if (this.data.negativePrompt && this.data.negativePrompt.trim()) {
        params.negative_prompt = this.data.negativePrompt;
      }
      
      // 先读取图片文件再创建任务
      wx.getFileSystemManager().readFile({
        filePath: this.data.tempFilePath,
        encoding: 'base64',
        success: (res) => {
          params.img_base64 = `data:${this.data.imageMimeType};base64,` + res.data;
          this.createTask(taskType, params);
        },
        fail: (err) => {
          this.handleError('图片读取失败: ' + err.errMsg);
        }
      });
      return; // 提前返回，避免重复调用createTask
    }

    this.createTask(taskType, params);
  },

  createTask: function(taskType, params) {
    // 根据任务类型调用不同的API
    if (taskType === 'text_to_video_online') {
      api.createVideoTaskWithPromptProcessing(params)
        .then(res => {
          if (res.code === 200 && res.data.job_id) {
            this.setData({ jobId: res.data.job_id });
            this.startPolling();
          } else {
            this.handleError('创建任务失败: ' + (res.message || ''));
          }
        })
        .catch((err) => {
         // 优先使用后端返回的错误信息
         const errorMessage = err.data && err.data.error ? err.data.error : (err.errMsg || '网络错误');
         this.handleError('创建任务请求失败: ' + errorMessage);
        });
    } else {
      // 其他类型使用原来的API
      api.createVideoTask(params)
        .then(res => {
          if (res.code === 200 && res.data.job_id) {
            this.setData({ jobId: res.data.job_id });
            this.startPolling();
          } else {
            this.handleError('创建任务失败: ' + (res.message || ''));
          }
        })
        .catch((err) => {
         // 优先使用后端返回的错误信息
         const errorMessage = err.data && err.data.error ? err.data.error : (err.errMsg || '网络错误');
         this.handleError('创建任务请求失败: ' + errorMessage);
        });
    }
  },

  startPolling: function () {
    // 清除之前的定时器
    if (this.data.pollingTimer) {
      clearTimeout(this.data.pollingTimer);
    }
    this.setData({ pollingCount: 0 }); // 重置轮询计数器

    const poll = () => {
      // 增加超时检查，24次 * 5秒 = 120秒
      if (this.data.pollingCount >= 24) {
        this.handleError('任务超时，请稍后再试');
        return;
      }
      this.setData({ pollingCount: this.data.pollingCount + 1 });

      api.queryTaskResult(this.data.jobId)
        .then(res => {
          if (res.code === 200 && res.data) {
            const { status, video_url, error_message } = res.data;
            if (status === 'SUCCEEDED') {
              this.setData({ resultVideoUrl: video_url, loading: false, pollingTimer: null });
              clearTimeout(this.data.pollingTimer);
            } else if (status === 'FAILED') {
              this.handleError(error_message || '视频生成失败');
              clearTimeout(this.data.pollingTimer);
              this.setData({ pollingTimer: null });
            } else if (status === 'RUNNING' || status === 'PENDING') {
              // 继续轮询
              const timer = setTimeout(poll, 5000);
              this.setData({ pollingTimer: timer });
            } else {
              this.handleError(`未知的任务状态: ${status}`);
              clearTimeout(this.data.pollingTimer);
              this.setData({ pollingTimer: null });
            }
          } else {
            this.handleError(res.message || '查询任务失败');
            clearTimeout(this.data.pollingTimer);
            this.setData({ pollingTimer: null });
          }
        })
        .catch((err) => {
          this.handleError('查询任务请求失败: ' + (err.errMsg || '网络错误'));
          clearTimeout(this.data.pollingTimer);
          this.setData({ pollingTimer: null });
        });
    };
    poll();
  },

  handleError: function (error) {
   if (this.data.pollingTimer) {
     clearTimeout(this.data.pollingTimer);
   }
   this.setData({ loading: false, pollingTimer: null });

   let errorMessage = '未知错误，请稍后再试';

   // 判断 error 是字符串还是对象
   if (typeof error === 'string') {
     errorMessage = error;
   } else if (error && error.data && error.data.error) {
     // 处理后端返回的错误 { "error": "..." }
     errorMessage = error.data.error;
   } else if (error && error.errMsg) {
     // 处理 wx.request 的 fail 回调错误
     errorMessage = error.errMsg;
   }

   // 使用Modal提示错误，方便查看详细信息
   wx.showModal({
     title: '生成失败',
     content: errorMessage,
     showCancel: false
   });
 },

  saveVideo: function () {
    if (!this.data.resultVideoUrl) return;
    wx.showLoading({ title: '保存中...' });
    wx.downloadFile({
      url: this.data.resultVideoUrl,
      success: (res) => {
        if (res.statusCode === 200) {
          // API返回的是GIF，所以使用saveImageToPhotosAlbum
          wx.saveImageToPhotosAlbum({
            filePath: res.tempFilePath,
            success: () => {
              wx.hideLoading();
              wx.showToast({ title: '已保存到相册' });
            },
            fail: (err) => {
              wx.hideLoading();
              // 检查是否是权限问题
              if (err.errMsg && err.errMsg.includes('auth deny')) {
                wx.showToast({ title: '保存失败，请授权相册访问', icon: 'none' });
              } else {
                wx.showToast({ title: '保存失败', icon: 'none' });
              }
            }
          });
        } else {
          wx.hideLoading();
          wx.showToast({ title: `下载失败: ${res.statusCode}`, icon: 'none' });
        }
      },
      fail: (err) => {
        wx.hideLoading();
        wx.showToast({ title: '下载失败: ' + err.errMsg, icon: 'none' });
      }
    });
  },

  // 转发GIF到好友
  shareToFriend: function() {
    if (!this.data.resultVideoUrl) return;
    
    wx.showLoading({ title: '准备分享...' });
    wx.downloadFile({
      url: this.data.resultVideoUrl,
      success: (res) => {
        if (res.statusCode === 200) {
          wx.hideLoading();
          // 预览图片，用户可以长按转发
          wx.previewImage({
            urls: [this.data.resultVideoUrl],
            current: this.data.resultVideoUrl,
            success: () => {
              console.log('预览成功');
            },
            fail: (err) => {
              wx.showToast({ title: '分享失败: ' + err.errMsg, icon: 'none' });
            }
          });
        } else {
          wx.hideLoading();
          wx.showToast({ title: `下载失败: ${res.statusCode}`, icon: 'none' });
        }
      },
      fail: (err) => {
        wx.hideLoading();
        wx.showToast({ title: '下载失败: ' + err.errMsg, icon: 'none' });
      }
    });
  },

  // 页面事件处理函数--监听页面分享到朋友圈
  onShareTimeline: function () {
    return {
      title: '表情包生成器',
      query: {}
    };
  },

  // 页面事件处理函数--监听页面分享给朋友
  onShareAppMessage: function () {
    return {
      title: '表情包生成器',
      path: '/pages/home/index'
    };
  }
});
