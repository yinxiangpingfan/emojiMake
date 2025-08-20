package com.example.emojimake.network

import android.content.Context
import android.os.Build
import coil.ImageLoader
import coil.decode.GifDecoder
import coil.decode.ImageDecoderDecoder

object CoilLoader {
    private var imageLoader: ImageLoader? = null

    fun getImageLoader(context: Context): ImageLoader {
        if (imageLoader == null) {
            // Reuse the unsafe OkHttpClient from ApiClient
            val unsafeOkHttpClient = ApiClient.getClient(context)
            
            imageLoader = ImageLoader.Builder(context)
                .okHttpClient(unsafeOkHttpClient)
                .components {
                    if (Build.VERSION.SDK_INT >= 28) {
                        add(ImageDecoderDecoder.Factory())
                    } else {
                        add(GifDecoder.Factory())
                    }
                }
                .build()
        }
        return imageLoader!!
    }
}
