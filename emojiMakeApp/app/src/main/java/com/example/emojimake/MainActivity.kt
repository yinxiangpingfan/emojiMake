package com.example.emojimake

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.example.emojimake.network.ApiClient
import com.example.emojimake.repository.UserPreferencesRepository
import com.example.emojimake.repository.UserRepository
import com.example.emojimake.repository.VideoRepository
import com.example.emojimake.ui.screens.*
import com.example.emojimake.ui.theme.EmojiMakeTheme
import com.example.emojimake.viewmodel.UserViewModel
import com.example.emojimake.viewmodel.VideoViewModel
import androidx.compose.runtime.LaunchedEffect

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)

        val apiService = ApiClient.create(this)
        val userPreferencesRepository = UserPreferencesRepository(this)
        val userRepository = UserRepository(apiService)
        val videoRepository = VideoRepository(apiService)

        val userViewModel: UserViewModel by lazy {
            ViewModelProvider(this, object : ViewModelProvider.Factory {
                override fun <T : ViewModel> create(modelClass: Class<T>): T {
                    return UserViewModel(userRepository, userPreferencesRepository) as T
                }
            })[UserViewModel::class.java]
        }

        val videoViewModel: VideoViewModel by lazy {
            ViewModelProvider(this, object : ViewModelProvider.Factory {
                override fun <T : ViewModel> create(modelClass: Class<T>): T {
                    return VideoViewModel(videoRepository) as T
                }
            })[VideoViewModel::class.java]
        }

        setContent {
            EmojiMakeTheme {
                val navController = rememberNavController()
                val authToken by userViewModel.authToken.collectAsState(initial = null)
                val userTokenExpired by userViewModel.tokenExpired.collectAsState()
                val videoTokenExpired by videoViewModel.tokenExpired.collectAsState()

                // 监听UserViewModel的token过期状态
                LaunchedEffect(userTokenExpired) {
                    if (userTokenExpired) {
                        // 清除用户认证信息
                        userViewModel.logout()
                        // 重置token过期状态
                        userViewModel.resetTokenExpired()
                        // 导航到登录界面
                        navController.navigate("auth") {
                            popUpTo("main") { inclusive = true }
                        }
                    }
                }
                
                // 监听VideoViewModel的token过期状态
                LaunchedEffect(videoTokenExpired) {
                    if (videoTokenExpired) {
                        // 清除用户认证信息
                        userViewModel.logout()
                        // 重置token过期状态
                        videoViewModel.resetTokenExpired()
                        // 导航到登录界面
                        navController.navigate("auth") {
                            popUpTo("main") { inclusive = true }
                        }
                    }
                }

                NavHost(
                    navController = navController,
                    startDestination = if (authToken != null) "main" else "auth"
                ) {
                    composable("auth") {
                        AuthScreen(userViewModel) {
                            navController.navigate("main") {
                                popUpTo("auth") { inclusive = true }
                            }
                        }
                    }
                    composable("main") {
                        MainScreen(
                            navController = navController,
                            onLogout = {
                                userViewModel.logout()
                                navController.navigate("auth") {
                                    popUpTo("main") { inclusive = true }
                                }
                            },
                            userViewModel = userViewModel
                        )
                    }
                    composable("change_password") {
                        ChangePasswordScreen(userViewModel, navController)
                    }
                    composable("text_to_video") {
                        TextToVideoScreen(videoViewModel)
                    }
                    composable("image_to_video") {
                        ImageToVideoScreen(videoViewModel)
                    }
                    composable("online_text_to_video") {
                        OnlineTextToVideoScreen(videoViewModel)
                    }
                }
            }
        }
    }
}