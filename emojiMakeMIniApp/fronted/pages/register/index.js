const api = require('../../utils/api.js');

Page({
  data: {
    phone: '',
    password: '',
    confirmPassword: '',
    showPassword: false,
    showConfirmPassword: false,
    agreedToTerms: false,
    canRegister: false
  },

  onLoad: function() {
    this.checkCanRegister();
  },

  // 切换密码显示状态
  togglePassword: function() {
    this.setData({
      showPassword: !this.data.showPassword
    });
  },

  // 切换确认密码显示状态
  toggleConfirmPassword: function() {
    this.setData({
      showConfirmPassword: !this.data.showConfirmPassword
    });
  },

  // 处理协议勾选
  onAgreementChange: function(e) {
    const agreed = e.detail.value.includes('agree');
    this.setData({
      agreedToTerms: agreed
    }, () => {
      this.checkCanRegister();
    });
  },

  // 检查是否可以注册
  checkCanRegister: function() {
    const { phone, password, confirmPassword, agreedToTerms } = this.data;
    const phoneValid = this.isValidPhone(phone);
    const passwordValid = password.length >= 8;
    const passwordMatch = password === confirmPassword && confirmPassword.length > 0;
    
    this.setData({
      canRegister: phoneValid && passwordValid && passwordMatch && agreedToTerms
    });
  },

  // 验证手机号格式
  isValidPhone: function(phone) {
    const phoneReg = /^1[3-9]\d{9}$/;
    return phoneReg.test(phone);
  },

  // 监听输入变化
  onPhoneInput: function(e) {
    this.setData({
      phone: e.detail.value
    }, () => {
      this.checkCanRegister();
    });
  },

  onPasswordInput: function(e) {
    this.setData({
      password: e.detail.value
    }, () => {
      this.checkCanRegister();
    });
  },

  onConfirmPasswordInput: function(e) {
    this.setData({
      confirmPassword: e.detail.value
    }, () => {
      this.checkCanRegister();
    });
  },

  // 注册功能
  register: function() {
    // 最终验证
    if (!this.data.canRegister) {
      wx.showToast({
        title: '请完善注册信息',
        icon: 'none'
      });
      return;
    }

    // 验证手机号
    if (!this.isValidPhone(this.data.phone)) {
      wx.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      });
      return;
    }

    // 验证密码长度
    if (this.data.password.length < 8) {
      wx.showToast({
        title: '密码至少需要8位',
        icon: 'none'
      });
      return;
    }

    // 验证密码一致性
    if (this.data.password !== this.data.confirmPassword) {
      wx.showToast({
        title: '两次输入的密码不一致',
        icon: 'none'
      });
      return;
    }

    // 验证协议勾选
    if (!this.data.agreedToTerms) {
      wx.showToast({
        title: '请先同意用户协议',
        icon: 'none'
      });
      return;
    }

    // 显示加载状态
    wx.showLoading({
      title: '注册中...',
      mask: true
    });

    // 调用注册API
    api.register(this.data.phone, this.data.password)
      .then(res => {
        wx.hideLoading();
        if (res.code === 0) {
          wx.showModal({
            title: '注册成功',
            content: '恭喜您注册成功！现在可以登录使用了。',
            showCancel: false,
            confirmText: '去登录',
            success: () => {
              wx.redirectTo({
                url: '/pages/login/index'
              });
            }
          });
        } else {
          wx.showToast({
            title: res.message || '注册失败',
            icon: 'none',
            duration: 3000
          });
        }
      })
      .catch(err => {
        wx.hideLoading();
        console.error('注册错误:', err);
        
        // 处理不同类型的错误
        let errorMessage = '注册失败，请稍后再试';
        if (err.data && err.data.message) {
          errorMessage = err.data.message;
        } else if (err.errMsg) {
          errorMessage = err.errMsg;
        }
        
        wx.showToast({
          title: errorMessage,
          icon: 'none',
          duration: 3000
        });
      });
  }
});