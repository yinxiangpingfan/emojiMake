package com.example.emojimake.network

import okhttp3.MultipartBody
import okhttp3.RequestBody
import retrofit2.http.*

interface ApiService {

    @Multipart
    @POST("/api/v1/users/register")
    suspend fun register(
        @Part("phone") phone: RequestBody,
        @Part("password") password: RequestBody
    ): BaseResponse<Unit>

    @Multipart
    @POST("/api/v1/users/login")
    suspend fun login(
        @Part("phone") phone: RequestBody,
        @Part("password") password: RequestBody
    ): BaseResponse<LoginResponseData>

    @Multipart
    @POST("/api/v1/users/change-password")
    suspend fun changePassword(
        @Part("newPassword") newPassword: RequestBody
    ): BaseResponse<Unit>

    @Multipart
    @POST("/api/v1/video/create")
    suspend fun createVideo(
        @Part("type") type: RequestBody,
        @Part("prompt") prompt: RequestBody,
        @Part("negative_prompt") negativePrompt: RequestBody?,
        @Part("size") size: RequestBody?,
        @Part("resolution") resolution: RequestBody?,
        @Part("img_base64") imgBase64: RequestBody?
    ): BaseResponse<CreateTaskResponseData>

    @Multipart
    @POST("/api/v1/video/create_with_prompt")
    suspend fun createVideoWithPrompt(
        @Part("role") role: RequestBody,
        @Part("source") source: RequestBody?,
        @Part("action") action: RequestBody,
        @Part("size") size: RequestBody
    ): BaseResponse<CreateTaskResponseData>

    @GET("/api/v1/video/query/{job_id}")
    suspend fun queryTask(
        @Path("job_id") jobId: String
    ): BaseResponse<QueryTaskResponseData>
}