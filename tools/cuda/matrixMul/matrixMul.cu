#include <stdio.h>

// 定义矩阵大小
#define N 32

// CUDA 核函数，用于矩阵乘法计算
__global__ void matrixMultiply(float* A, float* B, float* C)
{
    // 计算当前线程的全局索引
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;

    // 执行矩阵乘法运算
    float sum = 0.0;
    for (int k = 0; k < N; ++k) {
        sum += A[row * N + k] * B[k * N + col];
    }

    // 将结果保存到矩阵 C
    C[row * N + col] = sum;
}

int main()
{
    // 分配主机上的矩阵内存
    float* h_A = (float*)malloc(N * N * sizeof(float));
    float* h_B = (float*)malloc(N * N * sizeof(float));
    float* h_C = (float*)malloc(N * N * sizeof(float));

    // 初始化矩阵 A 和 B
    for (int i = 0; i < N * N; ++i) {
        h_A[i] = 1.0;
        h_B[i] = 2.0;
    }

    // 分配设备上的矩阵内存
    float* d_A, *d_B, *d_C;
    cudaMalloc(&d_A, N * N * sizeof(float));
    cudaMalloc(&d_B, N * N * sizeof(float));
    cudaMalloc(&d_C, N * N * sizeof(float));

    // 将矩阵 A 和 B 从主机内存复制到设备内存
    cudaMemcpy(d_A, h_A, N * N * sizeof(float), cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, N * N * sizeof(float), cudaMemcpyHostToDevice);

    // 定义 CUDA 核的网格和块大小
    dim3 gridSize(N / 16, N / 16);
    dim3 blockSize(16, 16);

    // 调用 CUDA 核函数进行矩阵乘法计算
    matrixMultiply<<<gridSize, blockSize>>>(d_A, d_B, d_C);

    // 将结果从设备内存复制到主机内存
    cudaMemcpy(h_C, d_C, N * N * sizeof(float), cudaMemcpyDeviceToHost);

    // 打印结果矩阵的一部分
    for (int i = 0; i < N; ++i) {
        for (int j = 0; j < N; ++j) {
            printf("%f ", h_C[i * N + j]);
        }
        printf("\n");
    }

    // 释放内存
    free(h_A);
    free(h_B);
    free(h_C);
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);

    return 0;
}
