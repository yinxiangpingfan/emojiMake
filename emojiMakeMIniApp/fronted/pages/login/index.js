const api = require('../../utils/api.js');

Page({
  data: {
    phone: '',
    password: '',
    showPassword: false
  },

  togglePassword: function() {
    this.setData({
      showPassword: !this.data.showPassword
    });
  },

  login: function() {
    api.login(this.data.phone, this.data.password)
      .then(res => {
        if (res.code === 0) {
          wx.setStorageSync('token', res.data.token);
          wx.showToast({
            title: '登录成功',
            icon: 'success'
          });
          wx.reLaunch({
            url: '/pages/home/index'
          });
        } else {
          wx.showToast({
            title: res.message,
            icon: 'none'
          });
        }
      })
      .catch(err => {
        wx.showToast({
          title: '登录失败',
          icon: 'none'
        });
      });
  }
});