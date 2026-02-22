import numpy as np
import matplotlib.pyplot as plt
from scipy.linalg import lu, qr, svd, eigh
import seaborn as sns

# ==================== 矩阵与向量 ====================

def demonstrate_matrix_operations():
    """演示矩阵基本运算"""
    print("=== 矩阵基本运算 ===")

    # 创建矩阵
    A = np.array([[1, 2, 3],
                  [4, 5, 6],
                  [7, 8, 9]])

    B = np.array([[9, 8, 7],
                  [6, 5, 4],
                  [3, 2, 1]])

    print(f"矩阵 A:\n{A}")
    print(f"矩阵 B:\n{B}")

    # 矩阵加法
    print(f"A + B:\n{A + B}")

    # 数乘
    scalar = 2
    print(f"{scalar} * A:\n{scalar * A}")

    # 矩阵乘法
    print(f"A × B:\n{np.dot(A, B)}")

    # 转置
    print(f"Aᵀ:\n{A.T}")

def demonstrate_vectors():
    """演示向量运算"""
    print("\n=== 向量运算 ===")

    # 创建向量
    v1 = np.array([1, 2, 3])
    v2 = np.array([4, 5, 6])

    print(f"向量 v1: {v1}")
    print(f"向量 v2: {v2}")

    # 向量加法
    print(f"v1 + v2: {v1 + v2}")

    # 数乘
    print(f"2 * v1: {2 * v1}")

    # 点积（内积）
    dot_product = np.dot(v1, v2)
    print(f"v1 · v2: {dot_product}")

    # 向量模（长度）
    norm_v1 = np.linalg.norm(v1)
    norm_v2 = np.linalg.norm(v2)
    print(f"||v1||: {norm_v1:.3f}")
    print(f"||v2||: {norm_v2:.3f}")

    # 夹角余弦
    cos_angle = dot_product / (norm_v1 * norm_v2)
    angle = np.arccos(cos_angle)
    print(f"夹角 (弧度): {angle:.3f}")
    print(f"夹角 (角度): {np.degrees(angle):.3f}")

def visualize_vector_operations():
    """可视化向量运算"""
    plt.figure(figsize=(12, 4))

    # 向量加法
    plt.subplot(131)
    v1 = np.array([2, 1])
    v2 = np.array([1, 3])
    v_sum = v1 + v2

    # 绘制向量
    plt.quiver(0, 0, v1[0], v1[1], angles='xy', scale_units='xy', scale=1,
               color='r', label='v1', width=0.01)
    plt.quiver(0, 0, v2[0], v2[1], angles='xy', scale_units='xy', scale=1,
               color='b', label='v2', width=0.01)
    plt.quiver(0, 0, v_sum[0], v_sum[1], angles='xy', scale_units='xy', scale=1,
               color='g', label='v1+v2', width=0.01)

    # 绘制平行四边形
    plt.plot([v1[0], v_sum[0]], [v1[1], v_sum[1]], 'k--', alpha=0.5)
    plt.plot([v2[0], v_sum[0]], [v2[1], v_sum[1]], 'k--', alpha=0.5)

    plt.xlim(-1, 5)
    plt.ylim(-1, 5)
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('向量加法')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.axis('equal')

    # 向量投影
    plt.subplot(132)
    v1 = np.array([3, 1])
    v2 = np.array([1, 2])

    # 投影长度
    proj_length = np.dot(v1, v2) / np.linalg.norm(v2)
    proj_vector = proj_length * v2 / np.linalg.norm(v2)

    plt.quiver(0, 0, v1[0], v1[1], angles='xy', scale_units='xy', scale=1,
               color='r', label='v1', width=0.01)
    plt.quiver(0, 0, v2[0], v2[1], angles='xy', scale_units='xy', scale=1,
               color='b', label='v2', width=0.01)
    plt.quiver(0, 0, proj_vector[0], proj_vector[1], angles='xy', scale_units='xy', scale=1,
               color='g', label='投影', width=0.01)

    # 绘制投影线
    plt.plot([v1[0], proj_vector[0]], [v1[1], proj_vector[1]], 'k--', alpha=0.5)

    plt.xlim(-1, 4)
    plt.ylim(-1, 3)
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('向量投影')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.axis('equal')

    # 3D向量
    plt.subplot(133, projection='3d')
    v1 = np.array([1, 2, 3])
    v2 = np.array([3, 1, 2])

    # 叉积
    cross_product = np.cross(v1, v2)

    plt.quiver(0, 0, 0, v1[0], v1[1], v1[2], color='r', label='v1', arrow_length_ratio=0.1)
    plt.quiver(0, 0, 0, v2[0], v2[1], v2[2], color='b', label='v2', arrow_length_ratio=0.1)
    plt.quiver(0, 0, 0, cross_product[0], cross_product[1], cross_product[2],
               color='g', label='v1×v2', arrow_length_ratio=0.1)

    plt.xlabel('x')
    plt.ylabel('y')
    plt.title('3D向量叉积')
    plt.legend()

    plt.tight_layout()
    plt.show()

# ==================== 线性方程组 ====================

def demonstrate_linear_systems():
    """演示线性方程组求解"""
    print("\n=== 线性方程组求解 ===")

    # 方程组：
    # 2x + y = 5
    # x - y = 1
    # 3x + 2y = 8

    A = np.array([[2, 1],
                  [1, -1]])
    b = np.array([5, 1])

    print(f"系数矩阵 A:\n{A}")
    print(f"常数向量 b: {b}")

    # 检查矩阵是否可逆
    det_A = np.linalg.det(A)
    print(f"det(A) = {det_A}")

    if abs(det_A) > 1e-10:
        # 逆矩阵法
        A_inv = np.linalg.inv(A)
        x_inv = np.dot(A_inv, b)
        print(f"\n逆矩阵法解: x = {x_inv}")

        # 高斯消元法
        x_solve = np.linalg.solve(A, b)
        print(f"高斯消元法解: x = {x_solve}")

        # 验证解
        print(f"验证: Ax = {np.dot(A, x_solve)}")
        print(f"原方程: b = {b}")
    else:
        print("矩阵不可逆，需要使用其他方法")

def visualize_linear_systems():
    """可视化线性方程组"""
    plt.figure(figsize=(10, 4))

    # 两条直线的交点
    plt.subplot(121)
    x = np.linspace(-2, 4, 100)

    # 2x + y = 5  => y = 5 - 2x
    # x - y = 1   => y = x - 1
    y1 = 5 - 2*x
    y2 = x - 1

    plt.plot(x, y1, 'r-', label='2x + y = 5', linewidth=2)
    plt.plot(x, y2, 'b-', label='x - y = 1', linewidth=2)

    # 交点
    intersection = np.array([2, 1])
    plt.plot(intersection[0], intersection[1], 'go', markersize=8, label='解 (2,1)')

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('线性方程组的几何意义')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.axis('equal')

    # 三条直线的情况
    plt.subplot(122)
    x = np.linspace(-1, 3, 100)

    y1 = 2 - x
    y2 = 0.5*x + 0.5
    y3 = 1.5*x - 1

    plt.plot(x, y1, 'r-', label='x + y = 2', linewidth=2)
    plt.plot(x, y2, 'b-', label='x - 2y = -1', linewidth=2)
    plt.plot(x, y3, 'g-', label='3x - 2y = 2', linewidth=2)

    # 交点
    intersection = np.array([1, 1])
    plt.plot(intersection[0], intersection[1], 'ko', markersize=8, label='解 (1,1)')

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('三方程组的解')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.axis('equal')

    plt.tight_layout()
    plt.show()

# ==================== 向量空间 ====================

def demonstrate_vector_spaces():
    """演示向量空间概念"""
    print("\n=== 向量空间概念 ===")

    # 线性无关性检查
    vectors = np.array([[1, 2, 3],
                       [2, 4, 6],
                       [1, 0, 1]])

    print("向量组:")
    for i, v in enumerate(vectors):
        print(f"v{i+1} = {v}")

    # 计算矩阵的秩
    rank = np.linalg.matrix_rank(vectors)
    print(f"\n向量组的秩: {rank}")

    if rank < len(vectors):
        print("向量组线性相关")
    else:
        print("向量组线性无关")

    # 基的变换
    demonstrate_basis_change()

def demonstrate_basis_change():
    """演示基的变换"""
    print("\n=== 基的变换 ===")

    # 标准基到新基的变换
    standard_basis = np.array([[1, 0, 0],
                               [0, 1, 0],
                               [0, 0, 1]])

    new_basis = np.array([[1, 1, 0],
                         [0, 1, 1],
                         [1, 0, 1]])

    print(f"标准基:\n{standard_basis}")
    print(f"新基:\n{new_basis}")

    # 变换矩阵
    transform_matrix = new_basis
    print(f"变换矩阵:\n{transform_matrix}")

    # 在标准基下的向量
    v_standard = np.array([2, 3, 1])
    print(f"\n标准基下的向量: {v_standard}")

    # 转换到新基下的坐标
    v_new = np.linalg.solve(transform_matrix, v_standard)
    print(f"新基下的坐标: {v_new}")

    # 验证
    v_reconstructed = np.dot(transform_matrix, v_new)
    print(f"重构向量: {v_reconstructed}")

# ==================== 特征值与特征向量 ====================

def demonstrate_eigenvalues():
    """演示特征值和特征向量"""
    print("\n=== 特征值与特征向量 ===")

    # 创建矩阵
    A = np.array([[3, 1],
                  [1, 3]])

    print(f"矩阵 A:\n{A}")

    # 计算特征值和特征向量
    eigenvalues, eigenvectors = np.linalg.eig(A)

    print(f"\n特征值: {eigenvalues}")
    print("特征向量:")
    for i, vec in enumerate(eigenvectors.T):
        print(f"λ{i+1} = {eigenvalues[i]:.3f} -> v{i+1} = {vec}")

    # 验证特征方程 Av = λv
    print("\n验证 Av = λv:")
    for i in range(len(eigenvalues)):
        Av = np.dot(A, eigenvectors[:, i])
        lambda_v = eigenvalues[i] * eigenvectors[:, i]
        print(f"v{i+1}: |Av - λv| = {np.linalg.norm(Av - lambda_v):.10f}")

    visualize_eigenvectors()

def visualize_eigenvectors():
    """可视化特征向量"""
    A = np.array([[3, 1],
                  [1, 3]])

    eigenvalues, eigenvectors = np.linalg.eig(A)

    plt.figure(figsize=(8, 6))

    # 绘制单位圆
    theta = np.linspace(0, 2*np.pi, 100)
    x_unit = np.cos(theta)
    y_unit = np.sin(theta)
    plt.plot(x_unit, y_unit, 'k--', alpha=0.3, label='单位圆')

    # 绘制特征向量
    colors = ['r', 'b']
    for i, (val, vec, color) in enumerate(zip(eigenvalues, eigenvectors.T, colors)):
        # 归一化特征向量
        vec_normalized = vec / np.linalg.norm(vec)

        # 绘制特征向量
        plt.quiver(0, 0, vec_normalized[0], vec_normalized[1],
                   angles='xy', scale_units='xy', scale=1,
                   color=color, label=f'特征向量 {i+1} (λ={val:.1f})',
                   width=0.01)

        # 绘制变换后的向量
        transformed = np.dot(A, vec_normalized)
        plt.quiver(0, 0, transformed[0], transformed[1],
                   angles='xy', scale_units='xy', scale=1,
                   color=color, alpha=0.5, linestyle='--',
                   width=0.01)

    plt.xlim(-2, 2)
    plt.ylim(-2, 2)
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('特征向量及其变换')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.axis('equal')
    plt.show()

# ==================== 矩阵分解 ====================

def demonstrate_matrix_decomposition():
    """演示矩阵分解"""
    print("\n=== 矩阵分解 ===")

    # 创建一个可逆矩阵
    A = np.array([[4, 3],
                  [6, 3]])

    print(f"原始矩阵 A:\n{A}")

    # LU分解
    P, L, U = lu(A)
    print(f"\nLU分解:")
    print(f"P (置换矩阵):\n{P}")
    print(f"L (下三角):\n{L}")
    print(f"U (上三角):\n{U}")
    print(f"验证 P×L×U:\n{np.dot(P, np.dot(L, U))}")

    # QR分解
    Q, R = qr(A)
    print(f"\nQR分解:")
    print(f"Q (正交矩阵):\n{Q}")
    print(f"R (上三角):\n{R}")
    print(f"验证 Q×R:\n{np.dot(Q, R)}")
    print(f"Qᵀ×Q = {np.dot(Q.T, Q)} (应该接近单位矩阵)")

    # SVD分解
    U, S, Vt = svd(A)
    print(f"\nSVD分解:")
    print(f"U:\n{U}")
    print(f"奇异值 S: {S}")
    print(f"Vᵀ:\n{Vt}")

    # 重构矩阵
    S_matrix = np.zeros(A.shape)
    S_matrix[:len(S), :len(S)] = np.diag(S)
    A_reconstructed = np.dot(U, np.dot(S_matrix, Vt))
    print(f"重构矩阵:\n{A_reconstructed}")

def visualize_svd():
    """可视化SVD分解"""
    # 创建一个2×3矩阵
    A = np.array([[3, 2, 2],
                  [2, 3, -2]])

    U, S, Vt = svd(A)

    plt.figure(figsize=(15, 5))

    # 原始矩阵的热力图
    plt.subplot(141)
    sns.heatmap(A, annot=True, cmap='coolwarm', center=0, square=True)
    plt.title('原始矩阵 A')

    # U矩阵
    plt.subplot(142)
    sns.heatmap(U, annot=True, cmap='coolwarm', center=0, square=True)
    plt.title('U矩阵')

    # 奇异值矩阵
    S_matrix = np.zeros(A.shape)
    S_matrix[:len(S), :len(S)] = np.diag(S)
    plt.subplot(143)
    sns.heatmap(S_matrix, annot=True, cmap='coolwarm', center=0, square=True)
    plt.title('奇异值矩阵 Σ')

    # Vt矩阵
    plt.subplot(144)
    sns.heatmap(Vt, annot=True, cmap='coolwarm', center=0, square=True)
    plt.title('Vᵀ矩阵')

    plt.tight_layout()
    plt.show()

    # 降维演示
    print(f"\n奇异值: {S}")
    print(f"能量保留比例:")
    for k in range(1, len(S)+1):
        energy = np.sum(S[:k]**2) / np.sum(S**2)
        print(f"前 {k} 个奇异值: {energy:.1%}")

# ==================== 线性变换 ====================

def demonstrate_linear_transformations():
    """演示线性变换"""
    print("\n=== 线性变换 ===")

    # 定义变换矩阵
    transformations = {
        '缩放': np.array([[2, 0],
                         [0, 0.5]]),
        '旋转': np.array([[np.cos(np.pi/4), -np.sin(np.pi/4)],
                         [np.sin(np.pi/4), np.cos(np.pi/4)]]),
        '反射': np.array([[1, 0],
                         [0, -1]]),
        '剪切': np.array([[1, 0.5],
                         [0, 1]])
    }

    visualize_transformations(transformations)

def visualize_transformations(transformations):
    """可视化各种线性变换"""
    # 创建原始图形（正方形）
    square = np.array([[-1, -1],
                       [1, -1],
                       [1, 1],
                       [-1, 1],
                       [-1, -1]])

    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    axes = axes.flatten()

    for idx, (name, matrix) in enumerate(transformations.items()):
        ax = axes[idx]

        # 应用变换
        transformed = np.dot(square, matrix.T)

        # 绘制原始图形
        ax.plot(square[:, 0], square[:, 1], 'b-', linewidth=2, label='原始')
        ax.fill(square[:, 0], square[:, 1], 'b', alpha=0.2)

        # 绘制变换后的图形
        ax.plot(transformed[:, 0], transformed[:, 1], 'r-', linewidth=2, label='变换后')
        ax.fill(transformed[:, 0], transformed[:, 1], 'r', alpha=0.2)

        # 绘制坐标轴
        ax.axhline(y=0, color='k', linestyle='-', alpha=0.3)
        ax.axvline(x=0, color='k', linestyle='-', alpha=0.3)

        # 设置图形属性
        ax.grid(True, alpha=0.3)
        ax.legend()
        ax.set_title(f'{name}变换')
        ax.set_xlabel('x')
        ax.set_ylabel('y')
        ax.axis('equal')

        # 显示变换矩阵
        ax.text(0.02, 0.98, f'矩阵:\n{matrix}', transform=ax.transAxes,
                verticalalignment='top', fontsize=8,
                bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.8))

    plt.tight_layout()
    plt.show()

# ==================== 主函数 ====================

def main():
    """主函数，运行所有演示"""
    print("线性代数学习演示程序")
    print("=" * 50)

    # 运行各个模块
    demonstrate_matrix_operations()
    demonstrate_vectors()
    visualize_vector_operations()
    demonstrate_linear_systems()
    visualize_linear_systems()
    demonstrate_vector_spaces()
    demonstrate_eigenvalues()
    demonstrate_matrix_decomposition()
    visualize_svd()
    demonstrate_linear_transformations()

    print("\n演示完成！")
    print("建议：")
    print("1. 尝试修改矩阵和向量的值，观察结果变化")
    print("2. 理解每种运算的几何意义")
    print("3. 练习手动计算，再用程序验证结果")

if __name__ == "__main__":
    main()