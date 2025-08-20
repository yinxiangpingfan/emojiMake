package com.example.emojimake.repository

import com.example.emojimake.network.ApiService
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody.Companion.toRequestBody

class UserRepository(private val apiService: ApiService) {

    suspend fun register(phone: String, password: String) =
        apiService.register(
            phone.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
            password.toRequestBody("multipart/form-data".toMediaTypeOrNull())
        )

    suspend fun login(phone: String, password: String) =
        apiService.login(
            phone.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
            password.toRequestBody("multipart/form-data".toMediaTypeOrNull())
        )

    suspend fun changePassword(newPassword: String) =
        apiService.changePassword(
            newPassword.toRequestBody("multipart/form-data".toMediaTypeOrNull())
        )
}