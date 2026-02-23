"""
概率论与数理统计进阶 - 代码示例
考研数学一专用

包含内容:
1. 三大抽样分布的蒙特卡罗验证
2. 最大似然估计 vs 矩估计对比
3. 假设检验功效曲线
4. 一元线性回归完整案例

依赖: numpy, matplotlib, scipy, sympy
"""

import numpy as np
import matplotlib.pyplot as plt
from scipy import stats
from scipy.optimize import minimize_scalar


def sampling_distributions_demo():
    """三大抽样分布的蒙特卡罗验证"""
    print("\n" + "=" * 60)
    print("【1. 三大抽样分布 - 蒙特卡罗验证】")
    print("=" * 60)

    np.random.seed(42)
    n_simulations = 10000
    sample_size = 10
    mu, sigma = 5.0, 2.0  # 总体参数: N(5, 4)

    # --- χ² 分布验证 ---
    # 定理: (n-1)S²/σ² ~ χ²(n-1)
    print("\n--- χ² 分布 ---")
    print(f"总体: N({mu}, {sigma**2}), 样本量 n={sample_size}")
    print(f"理论: (n-1)S²/σ² ~ χ²({sample_size - 1})")

    chi2_stats = []
    for _ in range(n_simulations):
        sample = np.random.normal(mu, sigma, sample_size)
        s2 = np.var(sample, ddof=1)  # 样本方差(无偏)
        chi2_val = (sample_size - 1) * s2 / sigma**2
        chi2_stats.append(chi2_val)

    chi2_stats = np.array(chi2_stats)
    theoretical_mean = sample_size - 1
    theoretical_var = 2 * (sample_size - 1)

    print(f"模拟均值: {chi2_stats.mean():.4f}, 理论均值: {theoretical_mean}")
    print(f"模拟方差: {chi2_stats.var():.4f}, 理论方差: {theoretical_var}")

    # --- t 分布验证 ---
    # 定理: (X̄-μ)/(S/√n) ~ t(n-1)
    print("\n--- t 分布 ---")
    print(f"理论: (X̄-μ)/(S/√n) ~ t({sample_size - 1})")

    t_stats = []
    for _ in range(n_simulations):
        sample = np.random.normal(mu, sigma, sample_size)
        x_bar = np.mean(sample)
        s = np.std(sample, ddof=1)
        t_val = (x_bar - mu) / (s / np.sqrt(sample_size))
        t_stats.append(t_val)

    t_stats = np.array(t_stats)
    print(f"模拟均值: {t_stats.mean():.4f}, 理论均值: 0")
    print(f"模拟方差: {t_stats.var():.4f}, 理论方差: {(sample_size-1)/(sample_size-3):.4f}")

    # --- F 分布验证 ---
    # 定理: (S1²/σ1²) / (S2²/σ2²) ~ F(n1-1, n2-1)
    print("\n--- F 分布 ---")
    n1, n2 = 10, 15
    sigma1, sigma2 = 2.0, 3.0
    print(f"两个总体: N(0,{sigma1**2}) 和 N(0,{sigma2**2})")
    print(f"样本量: n1={n1}, n2={n2}")
    print(f"理论: (S1²/σ1²)/(S2²/σ2²) ~ F({n1-1}, {n2-1})")

    f_stats = []
    for _ in range(n_simulations):
        s1 = np.random.normal(0, sigma1, n1)
        s2 = np.random.normal(0, sigma2, n2)
        var1 = np.var(s1, ddof=1)
        var2 = np.var(s2, ddof=1)
        f_val = (var1 / sigma1**2) / (var2 / sigma2**2)
        f_stats.append(f_val)

    f_stats = np.array(f_stats)
    f_theory_mean = (n2 - 1) / (n2 - 3) if n2 > 3 else float('inf')
    print(f"模拟均值: {f_stats.mean():.4f}, 理论均值: {f_theory_mean:.4f}")

    # 可视化
    fig, axes = plt.subplots(1, 3, figsize=(15, 4))
    fig.suptitle('三大抽样分布 - 蒙特卡罗验证', fontsize=14)

    # χ² 分布
    x_chi2 = np.linspace(0, 30, 200)
    axes[0].hist(chi2_stats, bins=50, density=True, alpha=0.7, label='模拟')
    axes[0].plot(x_chi2, stats.chi2.pdf(x_chi2, sample_size - 1), 'r-', lw=2, label='理论')
    axes[0].set_title(f'χ²({sample_size-1}) 分布')
    axes[0].legend()

    # t 分布
    x_t = np.linspace(-5, 5, 200)
    axes[1].hist(t_stats, bins=50, density=True, alpha=0.7, label='模拟')
    axes[1].plot(x_t, stats.t.pdf(x_t, sample_size - 1), 'r-', lw=2, label='理论')
    axes[1].plot(x_t, stats.norm.pdf(x_t), 'g--', lw=1, label='N(0,1)')
    axes[1].set_title(f't({sample_size-1}) 分布')
    axes[1].legend()

    # F 分布
    x_f = np.linspace(0, 5, 200)
    axes[2].hist(f_stats, bins=50, density=True, alpha=0.7, range=(0, 5), label='模拟')
    axes[2].plot(x_f, stats.f.pdf(x_f, n1 - 1, n2 - 1), 'r-', lw=2, label='理论')
    axes[2].set_title(f'F({n1-1}, {n2-1}) 分布')
    axes[2].legend()

    plt.tight_layout()
    plt.show()


def estimation_comparison_demo():
    """最大似然估计 vs 矩估计对比"""
    print("\n" + "=" * 60)
    print("【2. MLE vs 矩估计 对比】")
    print("=" * 60)

    np.random.seed(42)

    # --- 示例1: 均匀分布 U(0, θ) ---
    print("\n--- 均匀分布 U(0, θ) 参数估计 ---")
    theta_true = 5.0
    sample_sizes = [5, 10, 20, 50, 100]

    print(f"真实参数: θ = {theta_true}")
    print(f"\n{'样本量':>6} {'矩估计θ̂':>10} {'MLE θ̂':>10} {'矩估计偏差':>12} {'MLE偏差':>10}")
    print("-" * 55)

    mle_results = []
    mom_results = []

    for n in sample_sizes:
        # 多次模拟取平均
        mle_list, mom_list = [], []
        for _ in range(1000):
            sample = np.random.uniform(0, theta_true, n)
            # 矩估计: E(X) = θ/2, 所以 θ̂ = 2X̄
            mom_est = 2 * np.mean(sample)
            # MLE: θ̂ = max(X1, ..., Xn)
            mle_est = np.max(sample)
            mle_list.append(mle_est)
            mom_list.append(mom_est)

        mle_avg = np.mean(mle_list)
        mom_avg = np.mean(mom_list)
        mle_results.append(mle_list)
        mom_results.append(mom_list)

        print(f"{n:>6d} {mom_avg:>10.4f} {mle_avg:>10.4f} "
              f"{mom_avg - theta_true:>+12.4f} {mle_avg - theta_true:>+10.4f}")

    print("\n分析:")
    print("  - 矩估计是无偏的（E(2X̄) = θ）")
    print("  - MLE有偏（E(X(n)) = nθ/(n+1) < θ），但偏差随n增大趋于0")
    print("  - MLE的方差更小，更集中在真值附近")

    # --- 示例2: 正态分布方差估计 ---
    print("\n--- 正态分布 N(μ, σ²) 方差估计 ---")
    mu_true, sigma_true = 0, 2
    n = 20
    print(f"真实参数: μ={mu_true}, σ²={sigma_true**2}")
    print(f"样本量: n={n}")

    biased_vars, unbiased_vars = [], []
    for _ in range(10000):
        sample = np.random.normal(mu_true, sigma_true, n)
        biased_vars.append(np.var(sample))  # MLE: 除以n
        unbiased_vars.append(np.var(sample, ddof=1))  # 矩估计: 除以n-1

    print(f"  MLE方差估计(除以n)均值: {np.mean(biased_vars):.4f} (有偏)")
    print(f"  无偏方差估计(除以n-1)均值: {np.mean(unbiased_vars):.4f} (无偏)")
    print(f"  真实σ² = {sigma_true**2}")

    # 可视化: 均匀分布估计对比
    fig, axes = plt.subplots(1, 2, figsize=(12, 5))
    fig.suptitle('参数估计方法对比', fontsize=14)

    # 大样本下的分布对比
    axes[0].hist(mom_results[-1], bins=50, density=True, alpha=0.6, label='矩估计 2X̄')
    axes[0].hist(mle_results[-1], bins=50, density=True, alpha=0.6, label='MLE max(Xi)')
    axes[0].axvline(theta_true, color='r', ls='--', lw=2, label=f'θ={theta_true}')
    axes[0].set_title(f'U(0, {theta_true}) 估计分布 (n={sample_sizes[-1]})')
    axes[0].legend()

    # 正态方差估计对比
    axes[1].hist(biased_vars, bins=50, density=True, alpha=0.6, label='MLE (÷n, 有偏)')
    axes[1].hist(unbiased_vars, bins=50, density=True, alpha=0.6, label='无偏 (÷n-1)')
    axes[1].axvline(sigma_true**2, color='r', ls='--', lw=2, label=f'σ²={sigma_true**2}')
    axes[1].set_title(f'N({mu_true},{sigma_true**2}) 方差估计 (n={n})')
    axes[1].legend()

    plt.tight_layout()
    plt.show()


def hypothesis_testing_demo():
    """假设检验功效曲线"""
    print("\n" + "=" * 60)
    print("【3. 假设检验与功效曲线】")
    print("=" * 60)

    # --- 单样本 t 检验示例 ---
    print("\n--- 单样本 t 检验 ---")
    print("H0: μ = 50  vs  H1: μ ≠ 50")
    print("已知样本: n=25, X̄=52.3, S=8.5, α=0.05")

    n = 25
    x_bar = 52.3
    s = 8.5
    mu0 = 50
    alpha = 0.05

    # 计算t统计量
    t_stat = (x_bar - mu0) / (s / np.sqrt(n))
    # 临界值
    t_critical = stats.t.ppf(1 - alpha / 2, n - 1)
    # p值
    p_value = 2 * (1 - stats.t.cdf(abs(t_stat), n - 1))

    print(f"\nt 统计量 = {t_stat:.4f}")
    print(f"临界值 t_{alpha/2}({n-1}) = ±{t_critical:.4f}")
    print(f"p 值 = {p_value:.4f}")
    print(f"结论: {'拒绝H0' if abs(t_stat) > t_critical else '不拒绝H0'} (α={alpha})")

    # --- 功效曲线 (Power Curve) ---
    print("\n--- 功效曲线 ---")
    print("检验功效 = P(拒绝H0 | H0为假) = 1 - β")

    # 参数设置
    sigma_known = 8.5  # 假设σ已知简化计算
    mu_range = np.linspace(44, 56, 100)  # 真实μ的范围
    sample_sizes_power = [10, 25, 50, 100]

    fig, axes = plt.subplots(1, 2, figsize=(14, 5))
    fig.suptitle('假设检验', fontsize=14)

    # 左图: t分布与拒绝域
    x = np.linspace(-4, 4, 200)
    y = stats.t.pdf(x, n - 1)
    axes[0].plot(x, y, 'b-', lw=2)
    axes[0].fill_between(x, y, where=(x < -t_critical), alpha=0.3, color='red', label='拒绝域')
    axes[0].fill_between(x, y, where=(x > t_critical), alpha=0.3, color='red')
    axes[0].axvline(t_stat, color='green', ls='--', lw=2, label=f't={t_stat:.2f}')
    axes[0].set_title(f't检验 (df={n-1}, α={alpha})')
    axes[0].set_xlabel('t 值')
    axes[0].legend()

    # 右图: 功效曲线
    for ns in sample_sizes_power:
        z_alpha = stats.norm.ppf(1 - alpha / 2)
        power = []
        for mu_true in mu_range:
            # 功效 = P(|Z| > z_α/2 | μ=mu_true)
            ncp = (mu_true - mu0) / (sigma_known / np.sqrt(ns))
            p = 1 - (stats.norm.cdf(z_alpha - ncp) - stats.norm.cdf(-z_alpha - ncp))
            power.append(p)
        axes[1].plot(mu_range, power, lw=2, label=f'n={ns}')

    axes[1].axhline(alpha, color='gray', ls=':', label=f'α={alpha}')
    axes[1].axvline(mu0, color='gray', ls='--', alpha=0.5)
    axes[1].set_xlabel('真实 μ')
    axes[1].set_ylabel('功效 (1-β)')
    axes[1].set_title('功效曲线 (H0: μ=50)')
    axes[1].legend()
    axes[1].set_ylim(0, 1.05)

    plt.tight_layout()
    plt.show()

    # --- 两类错误总结 ---
    print("\n两类错误总结:")
    print("  ┌──────────┬──────────────┬──────────────┐")
    print("  │          │ H0 为真      │ H0 为假      │")
    print("  ├──────────┼──────────────┼──────────────┤")
    print("  │ 拒绝 H0  │ α错误(弃真) │ 正确(功效1-β)│")
    print("  │ 不拒绝H0 │ 正确(1-α)   │ β错误(取伪)  │")
    print("  └──────────┴──────────────┴──────────────┘")
    print("  - α增大 → β减小（此消彼长）")
    print("  - 增大n可同时减小α和β")


def linear_regression_demo():
    """一元线性回归完整案例"""
    print("\n" + "=" * 60)
    print("【4. 一元线性回归】")
    print("=" * 60)

    np.random.seed(42)

    # 生成数据: Y = 2.5 + 1.8X + ε, ε ~ N(0, 1)
    n = 30
    beta0_true, beta1_true = 2.5, 1.8
    sigma_true = 1.0

    X = np.random.uniform(1, 10, n)
    epsilon = np.random.normal(0, sigma_true, n)
    Y = beta0_true + beta1_true * X + epsilon

    print(f"真实模型: Y = {beta0_true} + {beta1_true}X + ε, ε~N(0,{sigma_true**2})")
    print(f"样本量: n={n}")

    # --- 最小二乘估计 ---
    x_bar = np.mean(X)
    y_bar = np.mean(Y)
    Sxy = np.sum((X - x_bar) * (Y - y_bar))
    Sxx = np.sum((X - x_bar) ** 2)
    Syy = np.sum((Y - y_bar) ** 2)

    beta1_hat = Sxy / Sxx
    beta0_hat = y_bar - beta1_hat * x_bar

    print(f"\n最小二乘估计:")
    print(f"  β̂₁ = Sxy/Sxx = {Sxy:.4f}/{Sxx:.4f} = {beta1_hat:.4f} (真值: {beta1_true})")
    print(f"  β̂₀ = Ȳ - β̂₁X̄ = {y_bar:.4f} - {beta1_hat:.4f}×{x_bar:.4f} = {beta0_hat:.4f} (真值: {beta0_true})")

    # --- 回归效果评估 ---
    Y_hat = beta0_hat + beta1_hat * X
    residuals = Y - Y_hat

    SSR = np.sum((Y_hat - y_bar) ** 2)  # 回归平方和
    SSE = np.sum(residuals ** 2)         # 残差平方和
    SST = np.sum((Y - y_bar) ** 2)       # 总平方和

    R_squared = SSR / SST
    sigma_hat_sq = SSE / (n - 2)  # 残差方差估计

    print(f"\n回归分析:")
    print(f"  SST (总变差) = {SST:.4f}")
    print(f"  SSR (回归)   = {SSR:.4f}")
    print(f"  SSE (残差)   = {SSE:.4f}")
    print(f"  SST = SSR + SSE: {SST:.4f} ≈ {SSR + SSE:.4f}")
    print(f"  R² = SSR/SST = {R_squared:.4f} (越接近1越好)")
    print(f"  σ̂² = SSE/(n-2) = {sigma_hat_sq:.4f}")

    # --- 显著性检验 (F检验) ---
    F_stat = (SSR / 1) / (SSE / (n - 2))
    F_critical = stats.f.ppf(0.95, 1, n - 2)
    p_value_f = 1 - stats.f.cdf(F_stat, 1, n - 2)

    print(f"\n显著性检验 (F检验):")
    print(f"  F = MSR/MSE = {F_stat:.4f}")
    print(f"  F临界值 F_0.05(1, {n-2}) = {F_critical:.4f}")
    print(f"  p值 = {p_value_f:.6f}")
    print(f"  结论: {'回归方程显著' if F_stat > F_critical else '回归方程不显著'}")

    # --- β1 的 t 检验 ---
    se_beta1 = np.sqrt(sigma_hat_sq / Sxx)
    t_stat_beta1 = beta1_hat / se_beta1
    p_value_t = 2 * (1 - stats.t.cdf(abs(t_stat_beta1), n - 2))

    print(f"\nβ₁的t检验 (H0: β₁=0):")
    print(f"  t = β̂₁/SE(β̂₁) = {beta1_hat:.4f}/{se_beta1:.4f} = {t_stat_beta1:.4f}")
    print(f"  p值 = {p_value_t:.6f}")

    # --- 可视化 ---
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    fig.suptitle('一元线性回归分析', fontsize=14)

    # 散点图+回归线
    x_line = np.linspace(X.min(), X.max(), 100)
    y_line = beta0_hat + beta1_hat * x_line
    axes[0, 0].scatter(X, Y, alpha=0.7, label='数据点')
    axes[0, 0].plot(x_line, y_line, 'r-', lw=2,
                     label=f'Ŷ = {beta0_hat:.2f} + {beta1_hat:.2f}X')
    axes[0, 0].plot(x_line, beta0_true + beta1_true * x_line, 'g--', lw=1,
                     label=f'真实: Y = {beta0_true} + {beta1_true}X')
    axes[0, 0].set_xlabel('X')
    axes[0, 0].set_ylabel('Y')
    axes[0, 0].set_title(f'回归拟合 (R²={R_squared:.4f})')
    axes[0, 0].legend()

    # 残差图
    axes[0, 1].scatter(Y_hat, residuals, alpha=0.7)
    axes[0, 1].axhline(0, color='r', ls='--')
    axes[0, 1].set_xlabel('预测值 Ŷ')
    axes[0, 1].set_ylabel('残差 e')
    axes[0, 1].set_title('残差图')

    # 残差正态性检验 (Q-Q图)
    stats.probplot(residuals, plot=axes[1, 0])
    axes[1, 0].set_title('残差Q-Q图')

    # 残差直方图
    axes[1, 1].hist(residuals, bins=10, density=True, alpha=0.7, label='残差')
    x_norm = np.linspace(residuals.min(), residuals.max(), 100)
    axes[1, 1].plot(x_norm, stats.norm.pdf(x_norm, 0, np.sqrt(sigma_hat_sq)),
                     'r-', lw=2, label='正态分布')
    axes[1, 1].set_xlabel('残差')
    axes[1, 1].set_title('残差分布')
    axes[1, 1].legend()

    plt.tight_layout()
    plt.show()

    print("\n408考点总结:")
    print("  - 最小二乘法: β̂₁ = Sxy/Sxx, β̂₀ = Ȳ - β̂₁X̄")
    print("  - 决定系数: R² = SSR/SST = 1 - SSE/SST")
    print("  - F检验与t检验等价(一元回归时 F = t²)")
    print("  - 残差分析: 检验独立性、等方差性、正态性")


def main():
    """主函数"""
    print("=" * 60)
    print("概率论与数理统计进阶 - 代码示例")
    print("考研数学一专用")
    print("=" * 60)

    # 设置中文字体（如果可用）
    try:
        plt.rcParams['font.sans-serif'] = ['SimHei', 'Microsoft YaHei', 'DejaVu Sans']
        plt.rcParams['axes.unicode_minus'] = False
    except Exception:
        pass

    sampling_distributions_demo()
    estimation_comparison_demo()
    hypothesis_testing_demo()
    linear_regression_demo()

    print("\n" + "=" * 60)
    print("所有概率统计进阶示例运行完成!")
    print("=" * 60)


if __name__ == "__main__":
    main()
