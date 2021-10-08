import axios from "axios"

// 封装请求
export const request = (config) => {

    // 1、创建axios实例
    const net = axios.create({
        baseURL: "http://192.168.162.128:9090",
        timeout: 5000,
    })

    // 2、定义拦截器
    net.interceptors.request.use(
        config => {
            return config
        },
        error => {

        }
    )

    // 3、发送真正的网络请求
    return net(config)
}