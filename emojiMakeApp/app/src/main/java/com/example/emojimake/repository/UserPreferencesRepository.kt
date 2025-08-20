package com.example.emojimake.repository

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "user_prefs")

class UserPreferencesRepository(private val context: Context) {

    private object PreferencesKeys {
        val AUTH_TOKEN = stringPreferencesKey("auth_token")
        val USER_PHONE = stringPreferencesKey("user_phone") // 添加手机号存储键
    }

    val authToken: Flow<String?> = context.dataStore.data
        .map { preferences ->
            preferences[PreferencesKeys.AUTH_TOKEN]
        }
        
    // 添加获取用户手机号的Flow
    val userPhone: Flow<String?> = context.dataStore.data
        .map { preferences ->
            preferences[PreferencesKeys.USER_PHONE]
        }

    suspend fun saveAuthToken(token: String) {
        context.dataStore.edit { preferences ->
            preferences[PreferencesKeys.AUTH_TOKEN] = token
        }
    }

    // 添加保存用户手机号的方法
    suspend fun saveUserPhone(phone: String) {
        context.dataStore.edit { preferences ->
            preferences[PreferencesKeys.USER_PHONE] = phone
        }
    }

    suspend fun clearAuthToken() {
        context.dataStore.edit { preferences ->
            preferences.remove(PreferencesKeys.AUTH_TOKEN)
        }
    }
    
    // 添加清除用户手机号的方法
    suspend fun clearUserPhone() {
        context.dataStore.edit { preferences ->
            preferences.remove(PreferencesKeys.USER_PHONE)
        }
    }
    
    // 添加同时清除认证令牌和手机号的方法
    suspend fun clearAllUserData() {
        context.dataStore.edit { preferences ->
            preferences.remove(PreferencesKeys.AUTH_TOKEN)
            preferences.remove(PreferencesKeys.USER_PHONE)
        }
    }
}