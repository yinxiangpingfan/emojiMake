package com.example.emojimake.ui.screens

import android.app.DownloadManager
import android.content.ContentValues
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Build
import android.os.Environment
import android.provider.MediaStore
import android.widget.Toast
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.content.FileProvider
import coil.compose.AsyncImage
import com.example.emojimake.network.ApiClient
import com.example.emojimake.network.CoilLoader
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.File
import java.io.FileOutputStream
import java.net.URL

@Composable
fun SuccessResult(url: String, onReset: () -> Unit) {
    val context = LocalContext.current
    val imageLoader = CoilLoader.getImageLoader(context)

    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
        modifier = Modifier.fillMaxSize()
    ) {
        Text("生成成功！", style = MaterialTheme.typography.headlineSmall)
        Spacer(modifier = Modifier.height(16.dp))
        AsyncImage(
            model = url,
            contentDescription = "生成的GIF",
            imageLoader = imageLoader,
            modifier = Modifier
                .size(256.dp)
                .clip(RoundedCornerShape(12.dp))
        )
        Spacer(modifier = Modifier.height(16.dp))

        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 32.dp),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            Button(onClick = { saveImageToGallery(context, url) }) {
                Text("保存到本地")
            }
            Button(onClick = { shareGif(context, url) }) {
                Text("分享")
            }
        }

        Spacer(modifier = Modifier.height(16.dp))
        Button(onClick = onReset) {
            Text("再来一次")
        }
    }
}

fun saveImageToGallery(context: Context, imageUrl: String) {
    val fileName = "emoji_${System.currentTimeMillis()}.gif"
    val mimeType = "image/gif"

    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
        val contentValues = ContentValues().apply {
            put(MediaStore.MediaColumns.DISPLAY_NAME, fileName)
            put(MediaStore.MediaColumns.MIME_TYPE, mimeType)
            put(MediaStore.MediaColumns.RELATIVE_PATH, "${Environment.DIRECTORY_PICTURES}/EmojiMake")
        }

        val resolver = context.contentResolver
        val uri = resolver.insert(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, contentValues)

        uri?.let {
            CoroutineScope(Dispatchers.IO).launch {
                try {
                    val client = ApiClient.getClient(context)
                    val request = okhttp3.Request.Builder().url(imageUrl).build()
                    val response = client.newCall(request).execute()
                    if (response.isSuccessful) {
                        response.body?.byteStream()?.use { inputStream ->
                            resolver.openOutputStream(it)?.use { outputStream ->
                                inputStream.copyTo(outputStream)
                            }
                        }
                        withContext(Dispatchers.Main) {
                            Toast.makeText(context, "保存成功", Toast.LENGTH_SHORT).show()
                        }
                    } else {
                        throw Exception("Download failed: ${response.message}")
                    }
                } catch (e: Exception) {
                    withContext(Dispatchers.Main) {
                        Toast.makeText(context, "保存失败: ${e.message}", Toast.LENGTH_SHORT).show()
                    }
                }
            }
        }
    } else {
        val request = DownloadManager.Request(Uri.parse(imageUrl)).apply {
            setAllowedNetworkTypes(DownloadManager.Request.NETWORK_WIFI or DownloadManager.Request.NETWORK_MOBILE)
            setTitle(fileName)
            setMimeType(mimeType)
            setNotificationVisibility(DownloadManager.Request.VISIBILITY_VISIBLE_NOTIFY_COMPLETED)
            setDestinationInExternalPublicDir(Environment.DIRECTORY_PICTURES, "EmojiMake/$fileName")
        }
        val downloadManager = context.getSystemService(Context.DOWNLOAD_SERVICE) as DownloadManager
        downloadManager.enqueue(request)
        Toast.makeText(context, "开始下载...", Toast.LENGTH_SHORT).show()
    }
}

fun shareGif(context: Context, imageUrl: String) {
    CoroutineScope(Dispatchers.IO).launch {
        try {
            val client = ApiClient.getClient(context)
            val request = okhttp3.Request.Builder().url(imageUrl).build()
            val response = client.newCall(request).execute()

            if (!response.isSuccessful) throw Exception("Download failed: ${response.message}")

            val file = File(context.cacheDir, "shared_gif.gif")
            FileOutputStream(file).use { outputStream ->
                response.body?.byteStream()?.copyTo(outputStream)
            }

            val uri = FileProvider.getUriForFile(context, "${context.packageName}.provider", file)

            withContext(Dispatchers.Main) {
                val intent = Intent(Intent.ACTION_SEND).apply {
                    type = "image/gif"
                    putExtra(Intent.EXTRA_STREAM, uri)
                    addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION)
                }
                context.startActivity(Intent.createChooser(intent, "分享GIF"))
            }
        } catch (e: Exception) {
            withContext(Dispatchers.Main) {
                Toast.makeText(context, "分享失败: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }
}
