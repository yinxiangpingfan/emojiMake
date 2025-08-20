package com.example.emojimake.network

data class BaseResponse<T>(
    val code: Int,
    val message: String?,
    val error: String?,
    val data: T?
) {
    fun getErrorMessage(): String {
        return message ?: error ?: "未知错误"
    }
}

data class LoginResponseData(
    val token: String
)

data class CreateTaskResponseData(
    val job_id: String
)

data class QueryTaskResponseData(
    val job_id: String,
    val status: String,
    val video_url: String?,
    val error_message: String?
)