"""
高等数学进阶 - Python 代码示例与可视化
适用于考研数学一学习

包含内容：
1. 符号曲线积分计算
2. 格林定理数值验证
3. 3D 向量场可视化
4. 傅里叶级数逼近
5. 级数收敛性测试
"""

import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import sympy as sp
from scipy import integrate
import seaborn as sns

# 设置中文显示
plt.rcParams['font.sans-serif'] = ['SimHei', 'Microsoft YaHei']
plt.rcParams['axes.unicode_minus'] = False

# 设置绘图风格
sns.set_style("whitegrid")


# ==================== 1. 符号曲线积分计算 ====================

def symbolic_line_integral():
    """使用 SymPy 进行符号曲线积分计算"""
    print("=" * 60)
    print("【示例1】符号曲线积分计算")
    print("=" * 60)
    
    # 定义符号变量
    t, x, y = sp.symbols('t x y')
    
    # 示例1：第一型曲线积分
    # 计算 ∫_L (x² + y²) ds，L 是从 (0,0) 到 (1,1) 的直线段
    print("\n1. 第一型曲线积分：∫_L (x² + y²) ds")
    print("   L: 从 (0,0) 到 (1,1) 的直线段")
    
    # 参数方程：x = t, y = t, 0 ≤ t ≤ 1
    x_param = t
    y_param = t
    
    # 计算 ds = √[(dx/dt)² + (dy/dt)²] dt
    dx_dt = sp.diff(x_param, t)
    dy_dt = sp.diff(y_param, t)
    ds = sp.sqrt(dx_dt**2 + dy_dt**2)
    
    # 被积函数
    f = x_param**2 + y_param**2
    
    # 计算积分
    integrand = f * ds
    result = sp.integrate(integrand, (t, 0, 1))
    
    print(f"   参数方程：x = {x_param}, y = {y_param}")
    print(f"   ds = {ds} dt")
    print(f"   被积表达式：{integrand}")
    print(f"   积分结果：{result} = {float(result):.6f}")
    
    # 示例2：第二型曲线积分
    print("\n2. 第二型曲线积分：∫_L y dx - x dy")
    print("   L: 单位圆从 (1,0) 逆时针到 (0,1)")
    
    # 参数方程：x = cos(t), y = sin(t), 0 ≤ t ≤ π/2
    x_param = sp.cos(t)
    y_param = sp.sin(t)
    
    dx = sp.diff(x_param, t) * sp.Symbol('dt')
    dy = sp.diff(y_param, t) * sp.Symbol('dt')
    
    # 计算 ∫ (y dx - x dy)
    integrand2 = y_param * sp.diff(x_param, t) - x_param * sp.diff(y_param, t)
    result2 = sp.integrate(integrand2, (t, 0, sp.pi/2))
    
    print(f"   参数方程：x = cos(t), y = sin(t)")
    print(f"   积分表达式：{integrand2}")
    print(f"   积分结果：{result2} = {float(result2):.6f}")
    
    # 示例3：空间曲线积分
    print("\n3. 空间曲线积分：∫_L xyz ds")
    print("   L: 螺旋线 x=cos(t), y=sin(t), z=t, 0≤t≤2π")
    
    z = sp.symbols('z')
    x_param = sp.cos(t)
    y_param = sp.sin(t)
    z_param = t
    
    # 计算 ds
    dx_dt = sp.diff(x_param, t)
    dy_dt = sp.diff(y_param, t)
    dz_dt = sp.diff(z_param, t)
    ds = sp.sqrt(dx_dt**2 + dy_dt**2 + dz_dt**2)
    
    # 被积函数
    f = x_param * y_param * z_param
    integrand3 = f * ds
    result3 = sp.integrate(integrand3, (t, 0, 2*sp.pi))
    
    print(f"   ds = {ds}")
    print(f"   积分结果：{result3} = {float(result3):.6f}")


# ==================== 2. 格林定理验证 ====================

def verify_green_theorem():
    """数值验证格林定理"""
    print("\n" + "=" * 60)
    print("【示例2】格林定理数值验证")
    print("=" * 60)
    print("格林公式：∮_L P dx + Q dy = ∬_D (∂Q/∂x - ∂P/∂y) dxdy")
    print("验证区域：单位圆 x² + y² = 1")
    print("向量场：P = -y/(x²+y²), Q = x/(x²+y²)")
    
    # 方法1：直接计算曲线积分（参数方程）
    def line_integral_parametric():
        # 单位圆参数方程：x = cos(t), y = sin(t)
        def integrand(t):
            x = np.cos(t)
            y = np.sin(t)
            dx_dt = -np.sin(t)
            dy_dt = np.cos(t)
            
            # P = -y/(x²+y²), Q = x/(x²+y²)
            # 在单位圆上 x² + y² = 1
            P = -y
            Q = x
            
            return P * dx_dt + Q * dy_dt
        
        result, error = integrate.quad(integrand, 0, 2*np.pi)
        return result
    
    line_result = line_integral_parametric()
    print(f"\n曲线积分（参数方程法）：{line_result:.6f}")
    
    # 方法2：计算二重积分
    # ∂Q/∂x - ∂P/∂y 对于这个特殊场在单位圆内部的计算
    # 注意：这个向量场在原点奇异，实际上旋度在非原点处为0
    # 但沿闭曲线积分不为0（不满足单连通）
    
    print(f"理论值：2π = {2*np.pi:.6f}")
    print(f"相对误差：{abs(line_result - 2*np.pi) / (2*np.pi) * 100:.4f}%")
    
    # 可视化向量场
    visualize_vector_field_2d()
    
    print("\n注意：此向量场在原点奇异，区域非单连通，")
    print("      因此曲线积分 ≠ 0（虽然在非原点处旋度为0）")


def visualize_vector_field_2d():
    """可视化2D向量场与闭曲线"""
    fig, ax = plt.subplots(figsize=(10, 10))
    
    # 创建网格
    x = np.linspace(-2, 2, 20)
    y = np.linspace(-2, 2, 20)
    X, Y = np.meshgrid(x, y)
    
    # 向量场 F = (-y, x) / (x² + y²)（在单位圆外也适用）
    # 避免原点奇异
    R2 = X**2 + Y**2
    R2[R2 < 0.01] = 0.01  # 避免除零
    
    U = -Y / R2
    V = X / R2
    
    # 绘制向量场
    ax.quiver(X, Y, U, V, alpha=0.6, color='blue', width=0.003)
    
    # 绘制单位圆
    theta = np.linspace(0, 2*np.pi, 100)
    circle_x = np.cos(theta)
    circle_y = np.sin(theta)
    ax.plot(circle_x, circle_y, 'r-', linewidth=2.5, label='积分路径 L')
    
    # 标记方向
    arrow_t = np.pi / 4
    ax.annotate('', xy=(np.cos(arrow_t + 0.1), np.sin(arrow_t + 0.1)),
                xytext=(np.cos(arrow_t), np.sin(arrow_t)),
                arrowprops=dict(arrowstyle='->', color='red', lw=2))
    
    ax.set_xlim(-2, 2)
    ax.set_ylim(-2, 2)
    ax.set_aspect('equal')
    ax.grid(True, alpha=0.3)
    ax.set_title('向量场 F = (-y, x)/(x²+y²) 与单位圆路径', fontsize=14, fontweight='bold')
    ax.set_xlabel('x', fontsize=12)
    ax.set_ylabel('y', fontsize=12)
    ax.legend(fontsize=11)
    
    plt.tight_layout()
    plt.show()


# ==================== 3. 3D向量场可视化 ====================

def visualize_3d_vector_field():
    """3D向量场可视化（梯度、散度、旋度演示）"""
    print("\n" + "=" * 60)
    print("【示例3】3D向量场可视化")
    print("=" * 60)
    
    # 示例向量场：F = (y, -x, z)（有旋度）
    print("\n向量场：F = (y, -x, z)")
    print("计算旋度：curl F = ∇ × F")
    
    # 使用 SymPy 计算旋度
    x, y, z = sp.symbols('x y z')
    P, Q, R = y, -x, z
    
    # 旋度
    curl_x = sp.diff(R, y) - sp.diff(Q, z)
    curl_y = sp.diff(P, z) - sp.diff(R, x)
    curl_z = sp.diff(Q, x) - sp.diff(P, y)
    
    print(f"curl F = ({curl_x}, {curl_y}, {curl_z})")
    
    # 散度
    div_F = sp.diff(P, x) + sp.diff(Q, y) + sp.diff(R, z)
    print(f"div F = {div_F}")
    
    # 绘制向量场
    fig = plt.figure(figsize=(14, 6))
    
    # 子图1：向量场 F = (y, -x, z)
    ax1 = fig.add_subplot(121, projection='3d')
    
    # 创建3D网格（稀疏以便可视化）
    x_vals = np.linspace(-1, 1, 6)
    y_vals = np.linspace(-1, 1, 6)
    z_vals = np.linspace(-1, 1, 6)
    X, Y, Z = np.meshgrid(x_vals, y_vals, z_vals)
    
    # 向量场分量
    U = Y
    V = -X
    W = Z
    
    # 绘制向量场
    ax1.quiver(X, Y, Z, U, V, W, length=0.2, normalize=True, 
               color='blue', alpha=0.7, arrow_length_ratio=0.3)
    
    ax1.set_xlabel('X', fontsize=11)
    ax1.set_ylabel('Y', fontsize=11)
    ax1.set_zlabel('Z', fontsize=11)
    ax1.set_title('向量场 F = (y, -x, z)\ncurl F = (0, 0, -2)', 
                  fontsize=12, fontweight='bold')
    
    # 子图2：梯度场示例 ∇u（保守场）
    ax2 = fig.add_subplot(122, projection='3d')
    
    # 标量场 u = x² + y² + z²，梯度 ∇u = (2x, 2y, 2z)
    U2 = 2 * X
    V2 = 2 * Y
    W2 = 2 * Z
    
    ax2.quiver(X, Y, Z, U2, V2, W2, length=0.15, normalize=True,
               color='red', alpha=0.7, arrow_length_ratio=0.3)
    
    ax2.set_xlabel('X', fontsize=11)
    ax2.set_ylabel('Y', fontsize=11)
    ax2.set_zlabel('Z', fontsize=11)
    ax2.set_title('梯度场 ∇u（u = x²+y²+z²）\ncurl(∇u) = 0（无旋）', 
                  fontsize=12, fontweight='bold')
    
    plt.tight_layout()
    plt.show()


# ==================== 4. 傅里叶级数展开与逼近 ====================

def fourier_series_approximation():
    """傅里叶级数逼近演示（方波与锯齿波）"""
    print("\n" + "=" * 60)
    print("【示例4】傅里叶级数逼近")
    print("=" * 60)
    
    # 示例1：方波函数
    print("\n1. 方波函数 f(x) = { -1, -π<x<0; 1, 0<x<π }")
    
    def square_wave(x):
        """方波函数"""
        return np.where(np.sin(x) >= 0, 1, -1)
    
    def fourier_square_wave(x, n_terms):
        """方波的傅里叶级数逼近（奇函数，只有正弦项）"""
        result = np.zeros_like(x)
        for n in range(1, n_terms + 1, 2):  # 只有奇数项
            result += (4 / (n * np.pi)) * np.sin(n * x)
        return result
    
    # 绘制不同项数的逼近
    x = np.linspace(-np.pi, np.pi, 1000)
    
    fig, axes = plt.subplots(2, 3, figsize=(15, 10))
    fig.suptitle('傅里叶级数逼近方波函数', fontsize=16, fontweight='bold')
    
    n_terms_list = [1, 3, 5, 10, 20, 50]
    
    for idx, (ax, n_terms) in enumerate(zip(axes.flat, n_terms_list)):
        # 原函数
        ax.plot(x, square_wave(x), 'k--', linewidth=1.5, label='原函数', alpha=0.7)
        
        # 傅里叶逼近
        approx = fourier_square_wave(x, n_terms)
        ax.plot(x, approx, 'r-', linewidth=2, label=f'{n_terms} 项逼近')
        
        ax.set_xlim(-np.pi, np.pi)
        ax.set_ylim(-1.5, 1.5)
        ax.grid(True, alpha=0.3)
        ax.legend(fontsize=9)
        ax.set_title(f'n = {n_terms}', fontsize=11)
        ax.set_xlabel('x', fontsize=10)
        ax.set_ylabel('f(x)', fontsize=10)
    
    plt.tight_layout()
    plt.show()
    
    # 示例2：锯齿波函数
    print("\n2. 锯齿波函数 f(x) = x (-π < x < π)")
    
    def sawtooth_wave(x):
        """锯齿波（在 [-π, π] 内）"""
        return x
    
    def fourier_sawtooth(x, n_terms):
        """锯齿波的傅里叶级数逼近"""
        result = np.zeros_like(x)
        for n in range(1, n_terms + 1):
            result += (2 * (-1)**(n+1) / n) * np.sin(n * x)
        return result
    
    # 绘制收敛动画帧
    fig, axes = plt.subplots(2, 3, figsize=(15, 10))
    fig.suptitle('傅里叶级数逼近锯齿波 f(x) = x', fontsize=16, fontweight='bold')
    
    x = np.linspace(-np.pi, np.pi, 1000)
    
    for idx, (ax, n_terms) in enumerate(zip(axes.flat, n_terms_list)):
        # 原函数
        ax.plot(x, sawtooth_wave(x), 'k--', linewidth=1.5, label='f(x) = x', alpha=0.7)
        
        # 傅里叶逼近
        approx = fourier_sawtooth(x, n_terms)
        ax.plot(x, approx, 'b-', linewidth=2, label=f'{n_terms} 项逼近')
        
        ax.set_xlim(-np.pi, np.pi)
        ax.set_ylim(-4, 4)
        ax.grid(True, alpha=0.3)
        ax.legend(fontsize=9)
        ax.set_title(f'n = {n_terms}', fontsize=11)
        ax.set_xlabel('x', fontsize=10)
        ax.set_ylabel('f(x)', fontsize=10)
    
    plt.tight_layout()
    plt.show()
    
    print("\n观察：")
    print("1. 在间断点处出现吉布斯现象（约9%的过冲）")
    print("2. 项数越多，逼近越精确")
    print("3. 方波收敛较慢（奇数项），锯齿波收敛较快")


# ==================== 5. 级数收敛性测试 ====================

def series_convergence_tests():
    """级数收敛性判别法演示"""
    print("\n" + "=" * 60)
    print("【示例5】级数收敛性测试")
    print("=" * 60)
    
    # 测试1：比值判别法（达朗贝尔判别法）
    print("\n1. 比值判别法测试")
    print("   级数：Σ n! / n^n")
    
    def ratio_test(n_max=50):
        """比值判别法：计算 lim(a_{n+1}/a_n)"""
        ratios = []
        for n in range(1, n_max):
            a_n = np.math.factorial(n) / (n ** n)
            a_n1 = np.math.factorial(n + 1) / ((n + 1) ** (n + 1))
            ratio = a_n1 / a_n
            ratios.append(ratio)
        return ratios
    
    ratios = ratio_test()
    limit_ratio = ratios[-1]
    
    print(f"   计算 lim(a_(n+1)/a_n) ≈ {limit_ratio:.6f}")
    print(f"   判别：ρ = {limit_ratio:.3f} < 1，级数收敛")
    
    # 绘制比值变化
    fig, axes = plt.subplots(2, 2, figsize=(14, 10))
    fig.suptitle('级数收敛性判别法演示', fontsize=16, fontweight='bold')
    
    # 子图1：比值判别法
    ax1 = axes[0, 0]
    ax1.plot(range(1, len(ratios) + 1), ratios, 'b-o', markersize=4)
    ax1.axhline(y=1, color='r', linestyle='--', linewidth=2, label='临界值 ρ=1')
    ax1.axhline(y=1/np.e, color='g', linestyle='--', linewidth=1.5, 
                label=f'极限值 ≈ {1/np.e:.3f}')
    ax1.set_xlabel('n', fontsize=11)
    ax1.set_ylabel('a_{n+1} / a_n', fontsize=11)
    ax1.set_title('比值判别法：Σ n!/n^n', fontsize=12, fontweight='bold')
    ax1.grid(True, alpha=0.3)
    ax1.legend(fontsize=10)
    
    # 测试2：根值判别法（柯西判别法）
    print("\n2. 根值判别法测试")
    print("   级数：Σ (2n/(3n+1))^n")
    
    def root_test(n_max=100):
        """根值判别法：计算 lim(a_n^(1/n))"""
        roots = []
        for n in range(1, n_max):
            a_n = (2 * n / (3 * n + 1)) ** n
            root = a_n ** (1 / n)
            roots.append(root)
        return roots
    
    roots = root_test()
    limit_root = roots[-1]
    
    print(f"   计算 lim(a_n^(1/n)) ≈ {limit_root:.6f}")
    print(f"   判别：ρ = {limit_root:.3f} < 1，级数收敛")
    
    # 子图2：根值判别法
    ax2 = axes[0, 1]
    ax2.plot(range(1, len(roots) + 1), roots, 'g-o', markersize=4)
    ax2.axhline(y=1, color='r', linestyle='--', linewidth=2, label='临界值 ρ=1')
    ax2.axhline(y=2/3, color='b', linestyle='--', linewidth=1.5, 
                label=f'极限值 = 2/3 ≈ {2/3:.3f}')
    ax2.set_xlabel('n', fontsize=11)
    ax2.set_ylabel('a_n^(1/n)', fontsize=11)
    ax2.set_title('根值判别法：Σ (2n/(3n+1))^n', fontsize=12, fontweight='bold')
    ax2.grid(True, alpha=0.3)
    ax2.legend(fontsize=10)
    
    # 测试3：交错级数（莱布尼茨判别法）
    print("\n3. 交错级数测试")
    print("   级数：Σ (-1)^(n-1) / n = ln(2)")
    
    def alternating_series(n_max=100):
        """交错调和级数的部分和"""
        partial_sums = []
        s = 0
        for n in range(1, n_max):
            s += (-1) ** (n - 1) / n
            partial_sums.append(s)
        return partial_sums
    
    partial_sums = alternating_series()
    theoretical_sum = np.log(2)
    
    print(f"   部分和 S_100 ≈ {partial_sums[-1]:.6f}")
    print(f"   理论值 ln(2) ≈ {theoretical_sum:.6f}")
    print(f"   误差：{abs(partial_sums[-1] - theoretical_sum):.6e}")
    
    # 子图3：交错级数收敛
    ax3 = axes[1, 0]
    ax3.plot(range(1, len(partial_sums) + 1), partial_sums, 'r-', linewidth=1.5)
    ax3.axhline(y=theoretical_sum, color='b', linestyle='--', linewidth=2, 
                label=f'ln(2) ≈ {theoretical_sum:.4f}')
    ax3.set_xlabel('n（项数）', fontsize=11)
    ax3.set_ylabel('部分和 S_n', fontsize=11)
    ax3.set_title('交错级数：Σ (-1)^(n-1)/n → ln(2)', fontsize=12, fontweight='bold')
    ax3.grid(True, alpha=0.3)
    ax3.legend(fontsize=10)
    
    # 测试4：p-级数
    print("\n4. p-级数测试：Σ 1/n^p")
    
    def p_series_partial_sum(p, n_max=1000):
        """p-级数的部分和"""
        n = np.arange(1, n_max + 1)
        return np.cumsum(1 / n ** p)
    
    # 子图4：不同p值的p-级数
    ax4 = axes[1, 1]
    
    p_values = [0.5, 1.0, 1.5, 2.0]
    colors = ['red', 'orange', 'green', 'blue']
    
    for p, color in zip(p_values, colors):
        partial_sums_p = p_series_partial_sum(p, n_max=500)
        n_vals = np.arange(1, len(partial_sums_p) + 1)
        
        if p <= 1:
            label = f'p={p} (发散)'
            linestyle = '--'
        else:
            label = f'p={p} (收敛)'
            linestyle = '-'
        
        ax4.plot(n_vals, partial_sums_p, color=color, linestyle=linestyle, 
                linewidth=2, label=label)
    
    ax4.set_xlabel('n（项数）', fontsize=11)
    ax4.set_ylabel('部分和 S_n', fontsize=11)
    ax4.set_title('p-级数：Σ 1/n^p（p>1收敛，p≤1发散）', fontsize=12, fontweight='bold')
    ax4.set_xlim(0, 500)
    ax4.set_ylim(0, 50)
    ax4.grid(True, alpha=0.3)
    ax4.legend(fontsize=10)
    
    plt.tight_layout()
    plt.show()
    
    print("\n总结：")
    print("  - 比值判别法：适用于含阶乘、指数的级数")
    print("  - 根值判别法：适用于含n次方的级数")
    print("  - 莱布尼茨判别法：适用于交错级数")
    print("  - p-级数：p>1收敛，p≤1发散")


# ==================== 主函数 ====================

def main():
    """主函数，运行所有演示"""
    print("\n" + "=" * 60)
    print("高等数学进阶 - Python 演示程序（考研数学一）")
    print("=" * 60)
    print("\n本程序包含以下演示：")
    print("1. 符号曲线积分计算（SymPy）")
    print("2. 格林定理数值验证")
    print("3. 3D向量场可视化（梯度、散度、旋度）")
    print("4. 傅里叶级数逼近（方波与锯齿波）")
    print("5. 级数收敛性测试（比值法、根值法、交错级数、p-级数）")
    print("\n开始运行...\n")
    
    # 运行各个模块
    symbolic_line_integral()
    verify_green_theorem()
    visualize_3d_vector_field()
    fourier_series_approximation()
    series_convergence_tests()
    
    print("\n" + "=" * 60)
    print("演示完成！")
    print("=" * 60)
    print("\n学习建议：")
    print("1. 仔细观察傅里叶级数的收敛过程（吉布斯现象）")
    print("2. 理解向量场的旋度与散度的几何意义")
    print("3. 掌握各种级数审敛法的适用条件")
    print("4. 通过可视化加深对抽象概念的理解")
    print("\n配合 calculus_advanced_guide.md 理论文档学习效果更佳！")


if __name__ == "__main__":
    main()
