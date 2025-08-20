package com.example.emojimake.viewmodel

import android.app.Application
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.example.emojimake.network.ApiClient
import com.example.emojimake.network.BaseResponse
import com.example.emojimake.repository.UserPreferencesRepository
import com.example.emojimake.repository.UserRepository
import com.google.gson.Gson
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import retrofit2.HttpException

sealed class AuthResult {
    object Idle : AuthResult()
    object Loading : AuthResult()
    object Success : AuthResult()
    data class Error(val message: String) : AuthResult()
}

class UserViewModelFactory(private val application: Application) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        if (modelClass.isAssignableFrom(UserViewModel::class.java)) {
            @Suppress("UNCHECKED_CAST")
            return UserViewModel(
                UserRepository(ApiClient.create(application)),
                UserPreferencesRepository(application)
            ) as T
        }
        throw IllegalArgumentException("Unknown ViewModel class")
    }
}

class UserViewModel(
    private val userRepository: UserRepository,
    private val userPreferencesRepository: UserPreferencesRepository
) : ViewModel() {

    private val _authResult = MutableStateFlow<AuthResult>(AuthResult.Idle)
    val authResult: StateFlow<AuthResult> = _authResult.asStateFlow()
    
    // 添加一个状态来通知token过期
    private val _tokenExpired = MutableStateFlow(false)
    val tokenExpired: StateFlow<Boolean> = _tokenExpired.asStateFlow()

    val authToken = userPreferencesRepository.authToken
    val userPhone = userPreferencesRepository.userPhone

    private fun handleAuthError(response: BaseResponse<*>) {
        val errorMessage = when (response.message) {
            "Invalid phone number format" -> "无效的手机号格式"
            "user with this phone number already exists" -> "该手机号已被注册"
            "invalid phone or password" -> "手机号或密码无效"
            "Invalid or expired JWT" -> {
                // 当token过期时，设置tokenExpired为true
                _tokenExpired.value = true
                "登录状态无效或已过期"
            }
            else -> response.message ?: "未知错误"
        }
        _authResult.value = AuthResult.Error(errorMessage)
    }

    fun register(phone: String, password: String) {
        viewModelScope.launch {
            _authResult.value = AuthResult.Loading
            try {
                val response = userRepository.register(phone, password)
                if (response.code == 0) {
                    _authResult.value = AuthResult.Success
                } else {
                    response.message?.let {
                        handleAuthError(response)
                    } ?: run {
                        _authResult.value = AuthResult.Error("未知错误")
                    }
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                try {
                    val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                    errorResponse.message?.let {
                        handleAuthError(errorResponse)
                    } ?: run {
                        _authResult.value = AuthResult.Error("未知错误")
                    }
                } catch (parseError: Exception) {
                    _authResult.value = AuthResult.Error("网络错误: ${e.code()}")
                }
            } catch (e: Exception) {
                _authResult.value = AuthResult.Error(e.message ?: "未知错误")
            }
        }
    }

    fun login(phone: String, password: String) {
        viewModelScope.launch {
            _authResult.value = AuthResult.Loading
            
            try {
                val response = userRepository.login(phone, password)
                if (response.code == 0 && response.data != null) {
                    userPreferencesRepository.saveAuthToken(response.data.token)
                    userPreferencesRepository.saveUserPhone(phone) // 保存手机号
                    _authResult.value = AuthResult.Success
                } else {
                    response.message?.let {
                        handleAuthError(response)
                    } ?: run {
                        _authResult.value = AuthResult.Error("未知错误")
                    }
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                try {
                    val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                    errorResponse.message?.let {
                        handleAuthError(errorResponse)
                    } ?: run {
                        _authResult.value = AuthResult.Error("网络错误: ${e.code()}")
                    }
                } catch (parseError: Exception) {
                    _authResult.value = AuthResult.Error("网络错误: ${e.code()}")
                }
            } catch (e: Exception) {
                _authResult.value = AuthResult.Error(e.message ?: "未知错误")
            }
        }
    }

    fun changePassword(newPassword: String) {
        viewModelScope.launch {
            _authResult.value = AuthResult.Loading
            try {
                val response = userRepository.changePassword(newPassword)
                if (response.code == 0) {
                    _authResult.value = AuthResult.Success
                } else {
                    response.message?.let {
                        handleAuthError(response)
                    } ?: run {
                        _authResult.value = AuthResult.Error("未知错误")
                    }
                }
            } catch (e: HttpException) {
                val errorBody = e.response()?.errorBody()?.string()
                try {
                    val errorResponse = Gson().fromJson(errorBody, BaseResponse::class.java)
                    errorResponse.message?.let {
                        handleAuthError(errorResponse)
                    } ?: run {
                        _authResult.value = AuthResult.Error("网络错误: ${e.code()}")
                    }
                } catch (parseError: Exception) {
                    _authResult.value = AuthResult.Error("网络错误: ${e.code()}")
                }
            } catch (e: Exception) {
                _authResult.value = AuthResult.Error(e.message ?: "未知错误")
            }
        }
    }

    fun logout() {
        viewModelScope.launch {
            userPreferencesRepository.clearAllUserData() // 清除所有用户数据包括手机号
        }
    }

    fun resetAuthResult() {
        _authResult.value = AuthResult.Idle
    }
    
    // 添加重置token过期状态的方法
    fun resetTokenExpired() {
        _tokenExpired.value = false
    }
}