"""
线性代数进阶 - Python 代码示例与可视化
适用于考研数学一学习

包含内容：
1. 二次型化标准形（配方法与正交变换法）
2. 特征值分解可视化（2D 线性变换）
3. 矩阵的四个基本子空间可视化
4. PCA 降维应用案例
"""

import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import sympy as sp
from scipy import linalg
import seaborn as sns

# 设置中文显示
plt.rcParams['font.sans-serif'] = ['SimHei', 'Microsoft YaHei']
plt.rcParams['axes.unicode_minus'] = False

# 设置绘图风格
sns.set_style("whitegrid")


# ==================== 1. 二次型化标准形 ====================

def quadratic_form_standard():
    """二次型化标准形演示（配方法与正交变换法对比）"""
    print("=" * 60)
    print("【示例1】二次型化标准形")
    print("=" * 60)
    
    # 定义二次型：f = 2x² + 5y² + 5z² + 4xy - 4xz - 8yz
    print("\n二次型：f = 2x² + 5y² + 5z² + 4xy - 4xz - 8yz")
    
    # 矩阵表示
    A = np.array([
        [ 2,  2, -2],
        [ 2,  5, -4],
        [-2, -4,  5]
    ])
    
    print("\n矩阵表示：f = x^T A x")
    print("A =")
    print(A)
    
    # 方法1：正交变换法（求特征值）
    print("\n" + "-" * 40)
    print("方法1：正交变换法")
    print("-" * 40)
    
    # 求特征值和特征向量
    eigenvalues, eigenvectors = linalg.eigh(A)
    
    print(f"\n特征值：{eigenvalues}")
    print(f"标准形：f = {eigenvalues[0]:.2f}y₁² + {eigenvalues[1]:.2f}y₂² + {eigenvalues[2]:.2f}y₃²")
    
    # 正交矩阵
    Q = eigenvectors
    print("\n正交矩阵 Q（列为单位特征向量）：")
    print(Q)
    
    # 验证正交性
    print("\nQ^T Q =")
    print(Q.T @ Q)
    
    # 验证对角化
    Lambda = Q.T @ A @ Q
    print("\nQ^T A Q =")
    print(np.diag(np.diag(Lambda)))  # 只显示对角元素
    
    # 判断正定性
    print("\n" + "-" * 40)
    print("正定性判断")
    print("-" * 40)
    
    if all(eigenvalues > 0):
        print("所有特征值 > 0，矩阵正定")
    elif all(eigenvalues < 0):
        print("所有特征值 < 0，矩阵负定")
    elif all(eigenvalues >= 0):
        print("所有特征值 ≥ 0，矩阵半正定")
    elif all(eigenvalues <= 0):
        print("所有特征值 ≤ 0，矩阵半负定")
    else:
        print("特征值有正有负，矩阵不定")
    
    # 惯性指数
    p = np.sum(eigenvalues > 0)  # 正惯性指数
    q = np.sum(eigenvalues < 0)  # 负惯性指数
    
    print(f"\n正惯性指数 p = {p}")
    print(f"负惯性指数 q = {q}")
    print(f"符号差 s = p - q = {p - q}")
    
    # 可视化二次曲面
    visualize_quadratic_surface(A, eigenvalues, eigenvectors)


def visualize_quadratic_surface(A, eigenvalues, eigenvectors):
    """可视化二次型对应的二次曲面"""
    fig = plt.figure(figsize=(16, 6))
    
    # 子图1：原坐标系下的二次曲面
    ax1 = fig.add_subplot(131, projection='3d')
    
    # 创建网格
    u = np.linspace(-1, 1, 30)
    v = np.linspace(-1, 1, 30)
    U, V = np.meshgrid(u, v)
    
    # 在原坐标系下计算曲面（设 f = 1）
    # 2x² + 5y² + 5z² + 4xy - 4xz - 8yz = 1
    # 这是一个椭球面，我们取 z 的一个解
    
    X = U
    Y = V
    # 简化：只显示 z = 0 平面附近的等高线
    Z = np.zeros_like(X)
    
    # 计算二次型的值
    vals = np.zeros_like(X)
    for i in range(X.shape[0]):
        for j in range(X.shape[1]):
            vec = np.array([X[i,j], Y[i,j], Z[i,j]])
            vals[i,j] = vec.T @ A @ vec
    
    surf1 = ax1.contour(X, Y, vals, levels=15, cmap='viridis')
    ax1.set_xlabel('x₁', fontsize=10)
    ax1.set_ylabel('x₂', fontsize=10)
    ax1.set_zlabel('f(x)', fontsize=10)
    ax1.set_title('原坐标系下的二次型', fontsize=12, fontweight='bold')
    
    # 子图2：标准形坐标系（特征向量坐标系）
    ax2 = fig.add_subplot(132, projection='3d')
    
    # 绘制特征向量
    origin = np.array([0, 0, 0])
    colors = ['r', 'g', 'b']
    labels = ['ξ₁', 'ξ₂', 'ξ₃']
    
    for i in range(3):
        vec = eigenvectors[:, i]
        ax2.quiver(origin[0], origin[1], origin[2],
                  vec[0], vec[1], vec[2],
                  color=colors[i], arrow_length_ratio=0.15, 
                  linewidth=2.5, label=labels[i])
    
    ax2.set_xlim([-1, 1])
    ax2.set_ylim([-1, 1])
    ax2.set_zlim([-1, 1])
    ax2.set_xlabel('x₁', fontsize=10)
    ax2.set_ylabel('x₂', fontsize=10)
    ax2.set_zlabel('x₃', fontsize=10)
    ax2.set_title('特征向量（主轴方向）', fontsize=12, fontweight='bold')
    ax2.legend(fontsize=10)
    
    # 子图3：特征值条形图
    ax3 = fig.add_subplot(133)
    
    bars = ax3.bar(range(1, 4), eigenvalues, color=['red', 'green', 'blue'], alpha=0.7)
    ax3.axhline(y=0, color='k', linestyle='--', linewidth=1)
    ax3.set_xlabel('特征值序号', fontsize=11)
    ax3.set_ylabel('特征值大小', fontsize=11)
    ax3.set_title(f'特征值（标准形系数）\n'
                  f'f = {eigenvalues[0]:.2f}y₁² + {eigenvalues[1]:.2f}y₂² + {eigenvalues[2]:.2f}y₃²', 
                  fontsize=12, fontweight='bold')
    ax3.set_xticks([1, 2, 3])
    ax3.set_xticklabels(['λ₁', 'λ₂', 'λ₃'])
    ax3.grid(True, alpha=0.3, axis='y')
    
    # 标注数值
    for i, (bar, val) in enumerate(zip(bars, eigenvalues)):
        height = bar.get_height()
        ax3.text(bar.get_x() + bar.get_width()/2, height + 0.1 if height > 0 else height - 0.3,
                f'{val:.2f}', ha='center', va='bottom' if height > 0 else 'top', 
                fontsize=11, fontweight='bold')
    
    plt.tight_layout()
    plt.show()


# ==================== 2. 特征值分解可视化（2D线性变换） ====================

def eigenvalue_decomposition_2d():
    """2D 线性变换与特征向量可视化"""
    print("\n" + "=" * 60)
    print("【示例2】特征值分解可视化（2D 线性变换）")
    print("=" * 60)
    
    # 定义一个2x2矩阵
    A = np.array([
        [3, 1],
        [0, 2]
    ])
    
    print("\n矩阵 A =")
    print(A)
    
    # 求特征值和特征向量
    eigenvalues, eigenvectors = linalg.eig(A)
    
    print(f"\n特征值：{eigenvalues}")
    print("\n特征向量：")
    for i, (val, vec) in enumerate(zip(eigenvalues, eigenvectors.T)):
        print(f"λ_{i+1} = {val:.4f}, ξ_{i+1} = {vec}")
    
    # 绘制变换效果
    fig, axes = plt.subplots(1, 3, figsize=(18, 6))
    
    # 子图1：单位圆变换
    ax1 = axes[0]
    
    # 单位圆上的点
    theta = np.linspace(0, 2*np.pi, 100)
    circle = np.array([np.cos(theta), np.sin(theta)])
    
    # 变换后
    transformed = A @ circle
    
    ax1.plot(circle[0], circle[1], 'b-', linewidth=2, label='变换前（单位圆）')
    ax1.plot(transformed[0], transformed[1], 'r-', linewidth=2, label='变换后（椭圆）')
    
    # 绘制特征向量
    for i in range(2):
        vec = eigenvectors[:, i].real
        val = eigenvalues[i].real
        
        # 变换前
        ax1.arrow(0, 0, vec[0], vec[1], head_width=0.1, head_length=0.1, 
                 fc='green', ec='green', linewidth=2, alpha=0.7)
        
        # 变换后（只伸缩，方向不变）
        transformed_vec = A @ vec
        ax1.arrow(0, 0, transformed_vec[0], transformed_vec[1], 
                 head_width=0.15, head_length=0.15,
                 fc='orange', ec='orange', linewidth=2, alpha=0.7)
    
    ax1.set_xlim(-4, 4)
    ax1.set_ylim(-4, 4)
    ax1.set_aspect('equal')
    ax1.grid(True, alpha=0.3)
    ax1.axhline(y=0, color='k', linewidth=0.5)
    ax1.axvline(x=0, color='k', linewidth=0.5)
    ax1.legend(fontsize=10)
    ax1.set_title('线性变换效果\n绿箭头→橙箭头（特征向量）', fontsize=12, fontweight='bold')
    ax1.set_xlabel('x', fontsize=11)
    ax1.set_ylabel('y', fontsize=11)
    
    # 子图2：网格变换
    ax2 = axes[1]
    
    # 创建网格
    x = np.linspace(-2, 2, 9)
    y = np.linspace(-2, 2, 9)
    
    # 水平线
    for yi in y:
        line_x = x
        line_y = np.full_like(x, yi)
        points = np.array([line_x, line_y])
        transformed_points = A @ points
        
        ax2.plot(line_x, line_y, 'b-', alpha=0.3, linewidth=1)
        ax2.plot(transformed_points[0], transformed_points[1], 'r-', alpha=0.6, linewidth=1.5)
    
    # 垂直线
    for xi in x:
        line_x = np.full_like(y, xi)
        line_y = y
        points = np.array([line_x, line_y])
        transformed_points = A @ points
        
        ax2.plot(line_x, line_y, 'b-', alpha=0.3, linewidth=1)
        ax2.plot(transformed_points[0], transformed_points[1], 'r-', alpha=0.6, linewidth=1.5)
    
    ax2.set_xlim(-8, 8)
    ax2.set_ylim(-5, 5)
    ax2.set_aspect('equal')
    ax2.grid(True, alpha=0.3)
    ax2.axhline(y=0, color='k', linewidth=0.5)
    ax2.axvline(x=0, color='k', linewidth=0.5)
    ax2.set_title('网格变换\n蓝色→红色', fontsize=12, fontweight='bold')
    ax2.set_xlabel('x', fontsize=11)
    ax2.set_ylabel('y', fontsize=11)
    
    # 子图3：特征向量方向
    ax3 = axes[2]
    
    # 绘制特征向量及其变换
    colors = ['green', 'purple']
    labels = ['ξ₁', 'ξ₂']
    
    for i, (val, vec, color, label) in enumerate(zip(eigenvalues, eigenvectors.T, colors, labels)):
        vec = vec.real
        val = val.real
        
        # 原向量
        ax3.arrow(0, 0, vec[0], vec[1], head_width=0.15, head_length=0.15,
                 fc=color, ec=color, linewidth=2.5, alpha=0.5, label=f'{label} (原)')
        
        # 变换后（伸缩 λ 倍）
        scaled_vec = val * vec
        ax3.arrow(0, 0, scaled_vec[0], scaled_vec[1], head_width=0.2, head_length=0.2,
                 fc=color, ec=color, linewidth=3, label=f'{label} × λ{i+1} = {val:.2f}')
        
        # 标注
        mid_point = scaled_vec / 2
        ax3.text(mid_point[0], mid_point[1], f'λ{i+1}={val:.2f}', 
                fontsize=10, fontweight='bold', color=color)
    
    ax3.set_xlim(-4, 4)
    ax3.set_ylim(-3, 3)
    ax3.set_aspect('equal')
    ax3.grid(True, alpha=0.3)
    ax3.axhline(y=0, color='k', linewidth=0.5)
    ax3.axvline(x=0, color='k', linewidth=0.5)
    ax3.legend(fontsize=9, loc='upper left')
    ax3.set_title('特征向量沿自身方向伸缩\nAξ = λξ', fontsize=12, fontweight='bold')
    ax3.set_xlabel('x', fontsize=11)
    ax3.set_ylabel('y', fontsize=11)
    
    plt.tight_layout()
    plt.show()


# ==================== 3. 矩阵的四个基本子空间 ====================

def four_fundamental_subspaces():
    """矩阵的四个基本子空间可视化"""
    print("\n" + "=" * 60)
    print("【示例3】矩阵的四个基本子空间")
    print("=" * 60)
    
    # 定义一个 3×4 矩阵
    A = np.array([
        [1,  2,  3,  4],
        [2,  4,  6,  8],
        [3,  5,  7,  9]
    ])
    
    print("\n矩阵 A (3×4) =")
    print(A)
    
    # 计算秩
    rank_A = np.linalg.matrix_rank(A)
    print(f"\nrank(A) = {rank_A}")
    
    # 1. 列空间 C(A)
    print("\n" + "-" * 40)
    print("1. 列空间 C(A)")
    print("-" * 40)
    print(f"维数：dim(C(A)) = rank(A) = {rank_A}")
    
    # 使用 SVD 找到列空间的正交基
    U, S, Vt = linalg.svd(A, full_matrices=True)
    col_space_basis = U[:, :rank_A]
    
    print("\n列空间的一组正交基：")
    print(col_space_basis)
    
    # 2. 零空间 N(A)
    print("\n" + "-" * 40)
    print("2. 零空间 N(A)（齐次方程 Ax = 0 的解空间）")
    print("-" * 40)
    
    null_dim = A.shape[1] - rank_A
    print(f"维数：dim(N(A)) = n - rank(A) = {A.shape[1]} - {rank_A} = {null_dim}")
    
    # 求零空间基础解系
    null_space_basis = linalg.null_space(A)
    
    print("\n零空间的一组基（基础解系）：")
    print(null_space_basis)
    
    # 验证 Ax = 0
    print("\n验证 A × (零空间向量) = 0：")
    for i in range(null_space_basis.shape[1]):
        vec = null_space_basis[:, i]
        result = A @ vec
        print(f"A × ξ{i+1} = {result}（应接近0向量）")
    
    # 3. 行空间 C(A^T)
    print("\n" + "-" * 40)
    print("3. 行空间 C(A^T)")
    print("-" * 40)
    print(f"维数：dim(C(A^T)) = rank(A) = {rank_A}")
    
    row_space_basis = Vt[:rank_A, :].T
    print("\n行空间的一组正交基：")
    print(row_space_basis)
    
    # 4. 左零空间 N(A^T)
    print("\n" + "-" * 40)
    print("4. 左零空间 N(A^T)（方程 A^T y = 0 的解空间）")
    print("-" * 40)
    
    left_null_dim = A.shape[0] - rank_A
    print(f"维数：dim(N(A^T)) = m - rank(A) = {A.shape[0]} - {rank_A} = {left_null_dim}")
    
    left_null_space_basis = linalg.null_space(A.T)
    print("\n左零空间的一组基：")
    print(left_null_space_basis)
    
    # 正交关系
    print("\n" + "-" * 40)
    print("正交关系验证")
    print("-" * 40)
    print("C(A) ⊥ N(A^T) 在 ℝ³ 中")
    print("C(A^T) ⊥ N(A) 在 ℝ⁴ 中")
    
    # 验证正交性
    if left_null_space_basis.shape[1] > 0 and col_space_basis.shape[1] > 0:
        orthogonality = col_space_basis.T @ left_null_space_basis
        print(f"\nC(A)基向量 · N(A^T)基向量 =")
        print(orthogonality)
        print("（应接近0矩阵）")
    
    # 可视化维数关系
    visualize_subspace_dimensions(A.shape[0], A.shape[1], rank_A)


def visualize_subspace_dimensions(m, n, r):
    """可视化四个基本子空间的维数关系"""
    fig, axes = plt.subplots(1, 2, figsize=(16, 6))
    
    # 子图1：维数柱状图
    ax1 = axes[0]
    
    subspaces = ['C(A)\n(列空间)', 'N(A)\n(零空间)', 
                 'C(A^T)\n(行空间)', 'N(A^T)\n(左零空间)']
    dimensions = [r, n - r, r, m - r]
    colors = ['skyblue', 'lightcoral', 'lightgreen', 'lightyellow']
    
    bars = ax1.bar(subspaces, dimensions, color=colors, alpha=0.8, edgecolor='black', linewidth=2)
    
    ax1.set_ylabel('维数', fontsize=12, fontweight='bold')
    ax1.set_title(f'四个基本子空间的维数\n矩阵 A: {m}×{n}, rank(A) = {r}', 
                  fontsize=14, fontweight='bold')
    ax1.grid(True, alpha=0.3, axis='y')
    
    # 标注数值
    for bar, dim in zip(bars, dimensions):
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width()/2, height + 0.1,
                f'{dim}', ha='center', va='bottom', fontsize=12, fontweight='bold')
    
    # 子图2：关系图
    ax2 = axes[1]
    ax2.axis('off')
    
    # 文本说明
    text = f"""
    矩阵 A: {m}×{n}, rank(A) = {r}
    
    四个基本子空间：
    
    1. 列空间 C(A) ⊆ ℝ^{m}
       - 维数：{r}
       - 由 A 的列向量张成
       - Ax 的所有可能结果
    
    2. 零空间 N(A) ⊆ ℝ^{n}
       - 维数：{n - r}
       - Ax = 0 的解空间
       - 基础解系
    
    3. 行空间 C(A^T) ⊆ ℝ^{n}
       - 维数：{r}
       - 由 A 的行向量张成
    
    4. 左零空间 N(A^T) ⊆ ℝ^{m}
       - 维数：{m - r}
       - A^T y = 0 的解空间
    
    正交关系：
       C(A) ⊥ N(A^T)  (在 ℝ^{m} 中)
       C(A^T) ⊥ N(A)  (在 ℝ^{n} 中)
    
    维数公式：
       dim(C(A)) + dim(N(A^T)) = {r} + {m-r} = {m}
       dim(C(A^T)) + dim(N(A)) = {r} + {n-r} = {n}
    """
    
    ax2.text(0.1, 0.5, text, fontsize=11, verticalalignment='center',
            family='monospace', bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))
    
    plt.tight_layout()
    plt.show()


# ==================== 4. PCA降维应用 ====================

def pca_application():
    """主成分分析（PCA）降维应用案例"""
    print("\n" + "=" * 60)
    print("【示例4】PCA 降维应用案例")
    print("=" * 60)
    
    # 生成2D数据（具有相关性）
    np.random.seed(42)
    n_samples = 200
    
    # 原始数据（在斜向分布）
    mean = [0, 0]
    cov = [[3, 2],    # 协方差矩阵（非对角，表示相关）
           [2, 2]]
    
    data = np.random.multivariate_normal(mean, cov, n_samples)
    
    print(f"\n生成 {n_samples} 个2维数据点")
    print(f"协方差矩阵：")
    print(cov)
    
    # 数据中心化
    data_centered = data - data.mean(axis=0)
    
    # 计算协方差矩阵
    cov_matrix = np.cov(data_centered.T)
    print(f"\n样本协方差矩阵：")
    print(cov_matrix)
    
    # 特征值分解（PCA的核心）
    eigenvalues, eigenvectors = linalg.eigh(cov_matrix)
    
    # 按特征值降序排列
    idx = eigenvalues.argsort()[::-1]
    eigenvalues = eigenvalues[idx]
    eigenvectors = eigenvectors[:, idx]
    
    print(f"\n特征值（方差）：{eigenvalues}")
    print(f"方差贡献率：{eigenvalues / eigenvalues.sum() * 100}%")
    
    print(f"\n主成分方向（特征向量）：")
    for i, vec in enumerate(eigenvectors.T):
        print(f"PC{i+1}: {vec}")
    
    # 投影到主成分
    data_pca = data_centered @ eigenvectors
    
    # 可视化
    visualize_pca(data, data_centered, data_pca, eigenvalues, eigenvectors)


def visualize_pca(data, data_centered, data_pca, eigenvalues, eigenvectors):
    """可视化 PCA 降维过程"""
    fig = plt.figure(figsize=(18, 6))
    
    # 子图1：原始数据与主成分方向
    ax1 = fig.add_subplot(131)
    
    ax1.scatter(data[:, 0], data[:, 1], alpha=0.6, s=30, c='blue', label='原始数据')
    
    # 绘制主成分方向
    mean = data.mean(axis=0)
    colors = ['red', 'green']
    labels = ['PC1（第一主成分）', 'PC2（第二主成分）']
    
    for i, (val, vec, color, label) in enumerate(zip(eigenvalues, eigenvectors.T, colors, labels)):
        # 特征向量缩放到特征值的平方根（标准差）
        scale = 2 * np.sqrt(val)
        ax1.arrow(mean[0], mean[1], vec[0] * scale, vec[1] * scale,
                 head_width=0.3, head_length=0.2, fc=color, ec=color, 
                 linewidth=3, label=label)
    
    ax1.set_xlabel('x₁', fontsize=11)
    ax1.set_ylabel('x₂', fontsize=11)
    ax1.set_title('原始数据与主成分方向\n红色：最大方差方向', fontsize=12, fontweight='bold')
    ax1.grid(True, alpha=0.3)
    ax1.legend(fontsize=9)
    ax1.set_aspect('equal')
    
    # 子图2：投影到第一主成分（降维）
    ax2 = fig.add_subplot(132)
    
    # 只保留第一主成分
    data_1d = data_pca[:, 0]
    
    # 在1D轴上显示
    ax2.scatter(data_1d, np.zeros_like(data_1d), alpha=0.6, s=30, c='red')
    ax2.set_xlabel('PC1（主成分1）', fontsize=11)
    ax2.set_yticks([])
    ax2.set_title(f'投影到第一主成分（降维到1D）\n保留方差：{eigenvalues[0]/eigenvalues.sum()*100:.1f}%', 
                  fontsize=12, fontweight='bold')
    ax2.grid(True, alpha=0.3, axis='x')
    ax2.set_ylim(-0.5, 0.5)
    
    # 子图3：PCA坐标系下的数据
    ax3 = fig.add_subplot(133)
    
    ax3.scatter(data_pca[:, 0], data_pca[:, 1], alpha=0.6, s=30, c='green')
    ax3.axhline(y=0, color='k', linewidth=0.5)
    ax3.axvline(x=0, color='k', linewidth=0.5)
    ax3.set_xlabel('PC1', fontsize=11)
    ax3.set_ylabel('PC2', fontsize=11)
    ax3.set_title('PCA坐标系下的数据\n（主成分不相关）', fontsize=12, fontweight='bold')
    ax3.grid(True, alpha=0.3)
    ax3.set_aspect('equal')
    
    plt.tight_layout()
    plt.show()
    
    # 绘制方差贡献率
    plot_variance_explained(eigenvalues)


def plot_variance_explained(eigenvalues):
    """绘制方差贡献率图"""
    fig, axes = plt.subplots(1, 2, figsize=(14, 5))
    
    # 方差比例
    variance_ratio = eigenvalues / eigenvalues.sum()
    cumulative_variance = np.cumsum(variance_ratio)
    
    # 子图1：各主成分方差贡献
    ax1 = axes[0]
    bars = ax1.bar(range(1, len(eigenvalues) + 1), variance_ratio * 100, 
                   color=['red', 'orange'], alpha=0.7, edgecolor='black', linewidth=2)
    ax1.set_xlabel('主成分', fontsize=11)
    ax1.set_ylabel('方差贡献率 (%)', fontsize=11)
    ax1.set_title('各主成分的方差贡献率', fontsize=12, fontweight='bold')
    ax1.set_xticks([1, 2])
    ax1.set_xticklabels(['PC1', 'PC2'])
    ax1.grid(True, alpha=0.3, axis='y')
    
    for bar, ratio in zip(bars, variance_ratio):
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width()/2, height + 1,
                f'{ratio*100:.1f}%', ha='center', va='bottom', 
                fontsize=11, fontweight='bold')
    
    # 子图2：累积方差贡献率
    ax2 = axes[1]
    ax2.plot(range(1, len(eigenvalues) + 1), cumulative_variance * 100, 
            'bo-', linewidth=2.5, markersize=10)
    ax2.axhline(y=95, color='r', linestyle='--', linewidth=2, label='95%阈值')
    ax2.set_xlabel('主成分数量', fontsize=11)
    ax2.set_ylabel('累积方差贡献率 (%)', fontsize=11)
    ax2.set_title('累积方差贡献率', fontsize=12, fontweight='bold')
    ax2.set_xticks([1, 2])
    ax2.set_xticklabels(['PC1', 'PC1+PC2'])
    ax2.set_ylim(0, 105)
    ax2.grid(True, alpha=0.3)
    ax2.legend(fontsize=10)
    
    for i, cum_var in enumerate(cumulative_variance):
        ax2.text(i + 1, cum_var * 100 + 2, f'{cum_var*100:.1f}%', 
                ha='center', fontsize=11, fontweight='bold')
    
    plt.tight_layout()
    plt.show()


# ==================== 主函数 ====================

def main():
    """主函数，运行所有演示"""
    print("\n" + "=" * 60)
    print("线性代数进阶 - Python 演示程序（考研数学一）")
    print("=" * 60)
    print("\n本程序包含以下演示：")
    print("1. 二次型化标准形（正交变换法）与正定性判断")
    print("2. 特征值分解可视化（2D 线性变换）")
    print("3. 矩阵的四个基本子空间")
    print("4. PCA 降维应用案例")
    print("\n开始运行...\n")
    
    # 运行各个模块
    quadratic_form_standard()
    eigenvalue_decomposition_2d()
    four_fundamental_subspaces()
    pca_application()
    
    print("\n" + "=" * 60)
    print("演示完成！")
    print("=" * 60)
    print("\n学习建议：")
    print("1. 理解二次型与实对称矩阵的关系")
    print("2. 掌握特征向量的几何意义（变换的主方向）")
    print("3. 记忆四个基本子空间的维数关系")
    print("4. 理解 PCA 的核心：协方差矩阵对角化")
    print("\n配合 linear_algebra_advanced_guide.md 理论文档学习效果更佳！")


if __name__ == "__main__":
    main()
