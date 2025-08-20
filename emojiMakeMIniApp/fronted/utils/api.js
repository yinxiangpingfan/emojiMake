const API_BASE_URL = 'https://82.156.59.17:8000';

const request = (url, method, data, requiresAuth = false, contentType = 'application/json') => {
  return new Promise((resolve, reject) => {
    const header = {
      'Content-Type': contentType
    };

    // 这是一个针对小程序 wx.request 的特殊处理。
    // 根据测试，对于需要授权的 multipart 请求，我们需要保留 Content-Type 头；
    // 而对于不需要授权的 multipart 请求（如登录、注册），则需要删除它，让小程序自动生成。
    // 对于 multipart/form-data，必须删除 Content-Type，
    // 让小程序自动生成包含 boundary 的请求头，否则后端无法解析。
    if (contentType === 'multipart/form-data') {
      delete header['Content-Type'];
      // 确保所有form-data的值都是字符串
      if (data) {
        for (const key in data) {
          if (Object.prototype.hasOwnProperty.call(data, key)) {
            data[key] = String(data[key]);
          }
        }
      }
    }

    if (requiresAuth) {
      const token = wx.getStorageSync('token');
      if (token) {
        header['Authorization'] = `Bearer ${token}`;
      } else {
        // 如果需要授权但没有token，直接拒绝请求
        reject('Authorization token not found');
        return;
      }
    }

    wx.request({
      url: `${API_BASE_URL}${url}`,
      method: method,
      data: data,
      header: header,
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data);
        } else {
          reject(res);
        }
      },
      fail: (err) => {
        reject(err);
      }
    });
  });
};

// 用户注册
const register = (phone, password) => {
  const url = '/api/v1/users/register';
  const data = {
    phone,
    password
  };
  return request(url, 'POST', data, false, 'multipart/form-data');
};

// 用户登录
const login = (phone, password) => {
  const url = '/api/v1/users/login';
  const data = {
    phone,
    password
  };
  return request(url, 'POST', data, false, 'multipart/form-data');
};

// 修改密码
const changePassword = (newPassword) => {
  const url = '/api/v1/users/change-password';
  const data = {
    newPassword
  };
  return request(url, 'POST', data, true, 'multipart/form-data'); // 此接口需要授权
};

const createVideoTask = (params) => {
  const url = '/api/v1/video/create';
  return request(url, 'POST', params, true, 'multipart/form-data'); // 此接口需要授权
};

// 新增的文生视频（联网）接口
const createVideoTaskWithPromptProcessing = (params) => {
  const url = '/api/v1/video/create_with_prompt';
  return request(url, 'POST', params, true, 'multipart/form-data'); // 此接口需要授权
};

const queryTaskResult = (jobId) => {
  const url = `/api/v1/video/query/${jobId}`;
  return request(url, 'GET', null, true); // 此接口需要授权
};

module.exports = {
  register,
  login,
  changePassword,
  createVideoTask,
  createVideoTaskWithPromptProcessing,
  queryTaskResult
};
