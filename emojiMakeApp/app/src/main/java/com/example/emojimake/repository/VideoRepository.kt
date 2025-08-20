package com.example.emojimake.repository

import com.example.emojimake.network.ApiService
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody.Companion.toRequestBody

class VideoRepository(private val apiService: ApiService) {

    suspend fun createVideo(
        type: String,
        prompt: String,
        negativePrompt: String?,
        size: String?,
        resolution: String?,
        imgBase64: String?
    ) = apiService.createVideo(
        type.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        prompt.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        negativePrompt?.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        size?.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        resolution?.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        imgBase64?.toRequestBody("multipart/form-data".toMediaTypeOrNull())
    )

    suspend fun createVideoWithPrompt(
        role: String,
        source: String?,
        action: String,
        size: String
    ) = apiService.createVideoWithPrompt(
        role.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        source?.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        action.toRequestBody("multipart/form-data".toMediaTypeOrNull()),
        size.toRequestBody("multipart/form-data".toMediaTypeOrNull())
    )

    suspend fun queryTask(jobId: String) = apiService.queryTask(jobId)
}