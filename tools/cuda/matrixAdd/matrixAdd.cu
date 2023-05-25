#include <stdio.h>

// 定义矩阵大小
#define N 32

// CUDA 核函数：矩阵加法
__global__ void matrixAdd(int *a, int *b, int *c)
{
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;

    if (row < N && col < N)
    {
        int index = row * N + col;
        c[index] = a[index] + b[index];
    }
}

int main()
{
    // 定义矩阵大小和字节数
    int numBytes = N * N * sizeof(int);

    // 分配主机内存
    int *h_a = (int *)malloc(numBytes);
    int *h_b = (int *)malloc(numBytes);
    int *h_c = (int *)malloc(numBytes);

    // 初始化矩阵数据
    for (int i = 0; i < N * N; i++)
    {
        h_a[i] = i;
        h_b[i] = i;
    }

    // 分配设备内存
    int *d_a, *d_b, *d_c;
    cudaMalloc((void **)&d_a, numBytes);
    cudaMalloc((void **)&d_b, numBytes);
    cudaMalloc((void **)&d_c, numBytes);

    // 将数据从主机内存复制到设备内存
    cudaMemcpy(d_a, h_a, numBytes, cudaMemcpyHostToDevice);
    cudaMemcpy(d_b, h_b, numBytes, cudaMemcpyHostToDevice);

    // 定义线程块和网格大小
    dim3 threadsPerBlock(16, 16);
    dim3 numBlocks(N / threadsPerBlock.x, N / threadsPerBlock.y);

    // 启动 CUDA 核函数进行矩阵加法运算
    matrixAdd<<<numBlocks, threadsPerBlock>>>(d_a, d_b, d_c);

    // 将结果从设备内存复制到主机内存
    cudaMemcpy(h_c, d_c, numBytes, cudaMemcpyDeviceToHost);

    // 打印结果
    for (int i = 0; i < N * N; i++)
    {
        printf("%d ", h_c[i]);
        if ((i + 1) % N == 0)
            printf("\n");
    }

    // 释放内存
    free(h_a);
    free(h_b);
    free(h_c);
    cudaFree(d_a);
    cudaFree(d_b);
    cudaFree(d_c);

    return 0;
}
