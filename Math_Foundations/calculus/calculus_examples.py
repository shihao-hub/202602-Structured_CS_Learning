import numpy as np
import matplotlib.pyplot as plt
from sympy import symbols, diff, integrate, limit, series, sin, cos, exp, log
import sympy as sp

# ==================== 极限与连续 ====================

def demonstrate_limits():
    """演示极限的计算"""
    x = symbols('x')

    # 计算各种极限
    examples = [
        ("lim(x→2) x²", limit(x**2, x, 2)),
        ("lim(x→0) sin(x)/x", limit(sin(x)/x, x, 0)),
        ("lim(x→∞) (1+1/x)^x", limit((1+1/x)**x, x, sp.oo)),
        ("lim(x→0) (e^x - 1)/x", limit((exp(x) - 1)/x, x, 0))
    ]

    print("=== 极限计算示例 ===")
    for desc, result in examples:
        print(f"{desc} = {result}")

    # 可视化极限
    visualize_limit()

def visualize_limit():
    """可视化函数 sin(x)/x 在 x=0 处的极限"""
    x_vals = np.linspace(-5, 5, 1000)
    y_vals = np.sin(x_vals) / x_vals
    y_vals[np.isnan(y_vals)] = 1  # 处理 x=0 的情况

    plt.figure(figsize=(10, 6))
    plt.plot(x_vals, y_vals, 'b-', label='sin(x)/x')
    plt.axhline(y=1, color='r', linestyle='--', label='极限值 = 1')
    plt.axvline(x=0, color='g', linestyle=':', alpha=0.5)
    plt.title('函数 sin(x)/x 在 x=0 处的极限')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.ylim(-0.5, 1.5)
    plt.show()

# ==================== 导数与微分 ====================

def demonstrate_derivatives():
    """演示导数的计算"""
    x = symbols('x')

    # 基本函数求导
    functions = [
        ("x³ + 2x² - x + 1", x**3 + 2*x**2 - x + 1),
        ("sin(x) + cos(x)", sin(x) + cos(x)),
        ("e^x * ln(x)", exp(x) * log(x)),
        ("x² * sin(x)", x**2 * sin(x))
    ]

    print("\n=== 导数计算示例 ===")
    for func_str, func in functions:
        derivative = diff(func, x)
        print(f"d/dx[{func_str}] = {derivative}")

    # 可视化导数
    visualize_derivatives()

def visualize_derivatives():
    """可视化函数及其导数"""
    x = np.linspace(-2, 2, 1000)
    f = x**3 - 3*x
    f_prime = 3*x**2 - 3

    plt.figure(figsize=(12, 8))

    plt.subplot(2, 2, 1)
    plt.plot(x, f, 'b-', linewidth=2, label='f(x) = x³ - 3x')
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('原函数')

    plt.subplot(2, 2, 2)
    plt.plot(x, f_prime, 'r-', linewidth=2, label="f'(x) = 3x² - 3")
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('导数')

    plt.subplot(2, 1, 2)
    plt.plot(x, f, 'b-', linewidth=2, label='f(x)')
    plt.plot(x, f_prime, 'r-', linewidth=2, label="f'(x)")
    plt.axhline(y=0, color='k', linestyle='-', alpha=0.3)
    plt.axvline(x=0, color='k', linestyle='-', alpha=0.3)
    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('函数与导数对比')

    plt.tight_layout()
    plt.show()

def demonstrate_critical_points():
    """演示极值点的求解"""
    x = symbols('x')
    f = x**3 - 6*x**2 + 9*x + 1
    f_prime = diff(f, x)
    f_second = diff(f_prime, x)

    # 求临界点
    critical_points = sp.solve(f_prime, x)

    print(f"\n=== 极值点分析 ===")
    print(f"函数: f(x) = {f}")
    print(f"一阶导数: f'(x) = {f_prime}")
    print(f"二阶导数: f''(x) = {f_second}")
    print(f"临界点: {critical_points}")

    # 判断极值类型
    for point in critical_points:
        second_value = f_second.subs(x, point)
        if second_value > 0:
            print(f"x = {point} 是极小值点 (f'' = {second_value})")
        elif second_value < 0:
            print(f"x = {point} 是极大值点 (f'' = {second_value})")
        else:
            print(f"x = {point} 需要进一步分析")

    # 可视化
    visualize_extrema()

def visualize_extrema():
    """可视化函数的极值点"""
    x = np.linspace(-1, 5, 1000)
    f = x**3 - 6*x**2 + 9*x + 1

    plt.figure(figsize=(10, 6))
    plt.plot(x, f, 'b-', linewidth=2, label='f(x) = x³ - 6x² + 9x + 1')

    # 标记极值点
    plt.plot([1, 3], [5, 3], 'ro', markersize=8, label='极值点')
    plt.annotate('极大值点 (1, 5)', xy=(1, 5), xytext=(1.5, 6),
                 arrowprops=dict(arrowstyle='->', color='red'))
    plt.annotate('极小值点 (3, 3)', xy=(3, 3), xytext=(3.5, 2),
                 arrowprops=dict(arrowstyle='->', color='red'))

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('函数的极值点')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.show()

# ==================== 积分 ====================

def demonstrate_integrals():
    """演示积分的计算"""
    x = symbols('x')

    # 不定积分
    print("\n=== 不定积分示例 ===")
    integrals = [
        ("∫x² dx", x**2),
        ("∫sin(x) dx", sin(x)),
        ("∫e^x dx", exp(x)),
        ("∫1/x dx", 1/x)
    ]

    for integral_str, func in integrals:
        result = integrate(func, x)
        print(f"{integral_str} = {result} + C")

    # 定积分
    print("\n=== 定积分示例 ===")
    definite_integrals = [
        ("∫[0,1] x² dx", (x**2, 0, 1)),
        ("∫[0,π] sin(x) dx", (sin(x), 0, sp.pi)),
        ("∫[0,1] e^x dx", (exp(x), 0, 1))
    ]

    for integral_str, (func, a, b) in definite_integrals:
        result = integrate(func, (x, a, b))
        print(f"{integral_str} = {result}")

def visualize_integration():
    """可视化定积分的几何意义"""
    x = np.linspace(0, 2, 1000)
    f = x**2

    plt.figure(figsize=(10, 6))

    # 填充积分区域
    x_fill = np.linspace(0, 1, 100)
    plt.fill_between(x_fill, 0, x_fill**2, alpha=0.3, color='blue',
                    label='积分区域 ∫[0,1] x² dx')

    # 绘制函数
    plt.plot(x, f, 'r-', linewidth=2, label='f(x) = x²')

    # 标记边界
    plt.axvline(x=0, color='k', linestyle='--', alpha=0.5)
    plt.axvline(x=1, color='k', linestyle='--', alpha=0.5)
    plt.axhline(y=0, color='k', linestyle='-', alpha=0.3)

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('定积分的几何意义：曲线下的面积')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.text(0.5, 0.1, '面积 = 1/3', fontsize=12, ha='center')
    plt.show()

# ==================== 级数 ====================

def demonstrate_series():
    """演示级数展开"""
    x = symbols('x')

    print("\n=== 泰勒级数展开示例 ===")

    # 不同函数的泰勒展开
    series_expansions = [
        ("e^x 在 x=0 处", exp(x), x, 0, 5),
        ("sin(x) 在 x=0 处", sin(x), x, 0, 6),
        ("cos(x) 在 x=0 处", cos(x), x, 0, 6),
        ("ln(1+x) 在 x=0 处", log(1+x), x, 0, 5)
    ]

    for desc, func, var, point, n_terms in series_expansions:
        expansion = series(func, var, point, n_terms)
        print(f"{desc}: {expansion}")

    visualize_series_convergence()

def visualize_series_convergence():
    """可视化级数收敛"""
    x = np.linspace(-2, 2, 1000)

    # e^x 的不同阶泰勒近似
    plt.figure(figsize=(12, 8))

    # 精确函数
    plt.plot(x, np.exp(x), 'k-', linewidth=2, label='e^x (精确)')

    # 不同阶数的泰勒近似
    colors = ['r', 'g', 'b', 'm', 'c']
    for n, color in enumerate([1, 2, 3, 4, 5], 1):
        # 计算泰勒级数
        x_sym = symbols('x')
        taylor_series = series(exp(x_sym), x_sym, 0, n+1).removeO()

        # 转换为数值函数
        taylor_func = sp.lambdify(x_sym, taylor_series, 'numpy')
        plt.plot(x, taylor_func(x), color=color, linestyle='--',
                label=f'{n}阶泰勒近似', alpha=0.8)

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('e^x 的泰勒级数逼近')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.ylim(-1, 8)
    plt.xlim(-2, 2)
    plt.show()

# ==================== 多元微积分 ====================

def demonstrate_multivariate():
    """演示多元函数微积分"""
    x, y = symbols('x y')

    print("\n=== 多元函数偏导数示例 ===")

    # 多元函数
    f = x**2 + y**2 + x*y
    fx = diff(f, x)
    fy = diff(f, y)
    fxx = diff(fx, x)
    fyy = diff(fy, y)
    fxy = diff(fx, y)

    print(f"函数: f(x,y) = {f}")
    print(f"∂f/∂x = {fx}")
    print(f"∂f/∂y = {fy}")
    print(f"∂²f/∂x² = {fxx}")
    print(f"∂²f/∂y² = {fyy}")
    print(f"∂²f/∂x∂y = {fxy}")

    # 梯度
    gradient = sp.Matrix([fx, fy])
    print(f"梯度 ∇f = {gradient}")

    visualize_multivariate_function()

def visualize_multivariate_function():
    """可视化多元函数"""
    fig = plt.figure(figsize=(15, 5))

    # 创建网格
    x = np.linspace(-2, 2, 100)
    y = np.linspace(-2, 2, 100)
    X, Y = np.meshgrid(x, y)

    # 函数 z = x² + y² + xy
    Z = X**2 + Y**2 + X*Y

    # 3D曲面图
    ax1 = fig.add_subplot(131, projection='3d')
    surf = ax1.plot_surface(X, Y, Z, cmap='viridis', alpha=0.8)
    ax1.set_xlabel('x')
    ax1.set_ylabel('y')
    ax1.set_zlabel('z')
    ax1.set_title('3D曲面图')

    # 等高线图
    ax2 = fig.add_subplot(132)
    contour = ax2.contour(X, Y, Z, levels=15, colors='black')
    ax2.clabel(contour, inline=True, fontsize=8)
    ax2.set_xlabel('x')
    ax2.set_ylabel('y')
    ax2.set_title('等高线图')
    ax2.set_aspect('equal')

    # 热力图
    ax3 = fig.add_subplot(133)
    im = ax3.imshow(Z, extent=[-2, 2, -2, 2], origin='lower', cmap='hot')
    ax3.set_xlabel('x')
    ax3.set_ylabel('y')
    ax3.set_title('热力图')
    plt.colorbar(im, ax=ax3)

    plt.tight_layout()
    plt.show()

# ==================== 常微分方程 ====================

def demonstrate_ode():
    """演示常微分方程"""
    print("\n=== 常微分方程示例 ===")

    # 简单的ODE示例
    t = symbols('t')
    y = sp.Function('y')

    # 一阶线性ODE: dy/dt + y = e^(-t)
    ode = sp.Eq(sp.diff(y(t), t) + y(t), sp.exp(-t))
    print(f"微分方程: {ode}")

    # 求解
    solution = sp.dsolve(ode)
    print(f"通解: {solution}")

    visualize_ode_solutions()

def visualize_ode_solutions():
    """可视化微分方程解"""
    # dy/dt = -y + 1 的解族
    t = np.linspace(0, 5, 1000)

    plt.figure(figsize=(10, 6))

    # 不同初始条件的解
    C_values = [-2, -1, 0, 1, 2, 3]
    colors = ['r', 'g', 'b', 'm', 'c', 'y']

    for C, color in zip(C_values, colors):
        # 解: y = 1 + C*e^(-t)
        y = 1 + C * np.exp(-t)
        plt.plot(t, y, color=color, linewidth=2, label=f'C = {C}')

    # 平衡解
    plt.axhline(y=1, color='k', linestyle='--', linewidth=2, label='平衡解 y = 1')

    plt.grid(True, alpha=0.3)
    plt.legend()
    plt.title('dy/dt = -y + 1 的解族')
    plt.xlabel('t')
    plt.ylabel('y(t)')
    plt.show()

# ==================== 主函数 ====================

def main():
    """主函数，运行所有演示"""
    print("高等数学学习演示程序")
    print("=" * 50)

    # 运行各个模块
    demonstrate_limits()
    demonstrate_derivatives()
    demonstrate_critical_points()
    demonstrate_integrals()
    visualize_integration()
    demonstrate_series()
    demonstrate_multivariate()
    demonstrate_ode()

    print("\n演示完成！")
    print("建议：")
    print("1. 仔细观察每个概念的可视化效果")
    print("2. 尝试修改代码中的参数，观察结果变化")
    print("3. 结合理论学习，加深理解")

if __name__ == "__main__":
    main()