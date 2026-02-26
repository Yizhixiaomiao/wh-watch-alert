// 移动端相关工具函数

/**
 * 检测是否为微信浏览器
 */
export const isWechatBrowser = () => {
    const ua = navigator.userAgent.toLowerCase()
    return ua.includes('micromessenger')
}

/**
 * 检测移动端平台
 */
export const detectMobilePlatform = () => {
    const ua = navigator.userAgent.toLowerCase()
    
    if (ua.includes('micromessenger')) {
        return 'wechat'
    }
    
    if (ua.includes('mobile') || ua.includes('android') || ua.includes('iphone')) {
        return 'mobile'
    }
    
    return 'web'
}

/**
 * 微信分享配置
 */
export const handleWechatShare = (shareData) => {
    if (!isWechatBrowser()) {
        return
    }

    // 检查是否有微信JS-SDK
    if (typeof wx === 'undefined') {
        console.warn('微信JS-SDK未加载')
        return
    }

    const defaultData = {
        title: '故障报修系统',
        desc: '快速提交故障报修，实时查看处理进度',
        link: window.location.href,
        imgUrl: '', // 分享图标URL
        success: () => {
        },
        cancel: () => {
        }
    }

    const finalData = { ...defaultData, ...shareData }

    // 配置微信分享
    wx.ready(() => {
        // 分享到朋友圈
        wx.updateTimelineShareData({
            title: finalData.title,
            link: finalData.link,
            imgUrl: finalData.imgUrl,
            success: finalData.success,
            cancel: finalData.cancel
        })

        // 分享给朋友
        wx.updateAppMessageShareData({
            title: finalData.title,
            desc: finalData.desc,
            link: finalData.link,
            imgUrl: finalData.imgUrl,
            success: finalData.success,
            cancel: finalData.cancel
        })
    })
}

/**
 * 加载微信JS-SDK
 */
export const loadWechatSDK = (appId, timestamp, nonceStr, signature) => {
    if (!isWechatBrowser()) {
        return Promise.resolve()
    }

    return new Promise((resolve, reject) => {
        // 动态加载微信JS-SDK
        const script = document.createElement('script')
        script.src = 'https://res.wx.qq.com/open/js/jweixin-1.6.0.js'
        script.onload = () => {
            wx.config({
                debug: false,
                appId: appId,
                timestamp: timestamp,
                nonceStr: nonceStr,
                signature: signature,
                jsApiList: [
                    'updateTimelineShareData',
                    'updateAppMessageShareData',
                    'onMenuShareTimeline',
                    'onMenuShareAppMessage'
                ]
            })

            wx.ready(() => {
                resolve()
            })

            wx.error((res) => {
                console.error('微信JS-SDK配置失败:', res)
                reject(res)
            })
        }
        script.onerror = reject
        document.head.appendChild(script)
    })
}

/**
 * 获取设备信息
 */
export const getDeviceInfo = () => {
    const ua = navigator.userAgent
    const device = {
        userAgent: ua,
        platform: navigator.platform,
        language: navigator.language,
        isWechat: isWechatBrowser(),
        isIOS: /iPad|iPhone|iPod/.test(ua),
        isAndroid: /Android/.test(ua),
        isMobile: /Mobile|Android|iPhone|iPad|iPod/.test(ua)
    }

    // 获取微信版本
    if (device.isWechat) {
        const match = ua.match(/MicroMessenger\/(\d+\.\d+\.\d+)/)
        device.wechatVersion = match ? match[1] : 'unknown'
    }

    // 获取屏幕信息
    device.screenWidth = window.screen.width
    device.screenHeight = window.screen.height
    device.viewportWidth = window.innerWidth
    device.viewportHeight = window.innerHeight

    return device
}

/**
 * 格式化时间显示
 */
export const formatTime = (timestamp) => {
    if (!timestamp) return ''
    
    const date = new Date(timestamp * 1000)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    
    // 小于1分钟
    if (diff < 60000) {
        return '刚刚'
    }
    
    // 小于1小时
    if (diff < 3600000) {
        return Math.floor(diff / 60000) + '分钟前'
    }
    
    // 小于1天
    if (diff < 86400000) {
        return Math.floor(diff / 3600000) + '小时前'
    }
    
    // 小于7天
    if (diff < 604800000) {
        return Math.floor(diff / 86400000) + '天前'
    }
    
    // 超过7天显示具体日期
    return date.toLocaleDateString()
}

/**
 * 复制到剪贴板
 */
export const copyToClipboard = (text) => {
    if (navigator.clipboard) {
        return navigator.clipboard.writeText(text)
    } else {
        // 兼容处理
        const textArea = document.createElement('textarea')
        textArea.value = text
        textArea.style.position = 'fixed'
        textArea.style.opacity = '0'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        
        return new Promise((resolve, reject) => {
            try {
                const successful = document.execCommand('copy')
                document.body.removeChild(textArea)
                successful ? resolve() : reject()
            } catch (err) {
                document.body.removeChild(textArea)
                reject(err)
            }
        })
    }
}

/**
 * 显示toast提示
 */
export const showToast = (message, type = 'info') => {
    if (typeof window !== 'undefined' && window.WeixinJSBridge) {
        // 微信环境使用微信toast
        WeixinJSBridge.invoke('showToast', {
            title: message,
            icon: type === 'success' ? 'success' : 'none',
            duration: 2000
        })
    } else {
        // 其他环境使用console或实现自定义toast
    }
}

/**
 * 拨打电话
 */
export const makePhoneCall = (phoneNumber) => {
    if (isWechatBrowser() && typeof wx !== 'undefined') {
        // 微信环境使用微信API
        wx.makePhoneCall({
            phoneNumber: phoneNumber
        })
    } else {
        // 普通环境
        window.location.href = `tel:${phoneNumber}`
    }
}

/**
 * 获取当前位置
 */
export const getCurrentPosition = () => {
    return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
            reject(new Error('浏览器不支持地理定位'))
            return
        }

        navigator.geolocation.getCurrentPosition(
            (position) => {
                resolve({
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude,
                    accuracy: position.coords.accuracy
                })
            },
            (error) => {
                reject(error)
            },
            {
                enableHighAccuracy: true,
                timeout: 10000,
                maximumAge: 60000
            }
        )
    })
}