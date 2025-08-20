const api = require('../../utils/api.js');

Page({
  data: {
    newPassword: ''
  },

  changePassword: function() {
    api.changePassword(this.data.newPassword)
      .then(res => {
        if (res.code === 0) {
          wx.showToast({
            title: '修改成功',
            icon: 'success'
          });
          wx.removeStorageSync('token');
          wx.navigateTo({
            url: '/pages/login/index'
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
          title: '修改失败',
          icon: 'none'
        });
      });
  }
});