import numpy as np
import matplotlib.pyplot as plt
from scipy import stats
import seaborn as sns
from collections import Counter

# ==================== 概率基础 ====================

def demonstrate_basic_probability():
    """演示基本概率概念"""
    print("=== 基本概率概念 ===")

    # 掷骰子示例
    print("\n掷骰子示例：")
    outcomes = list(range(1, 7))  # 骰子的6个面
    print(f"样本空间: {outcomes}")

    # 计算各种事件的概率
    # 事件A：掷出偶数
    event_A = [2, 4, 6]
    prob_A = len(event_A) / len(outcomes)
    print(f"事件A(偶数): {event_A}, P(A) = {prob_A}")

    # 事件B：掷出大于4的点数
    event_B = [5, 6]
    prob_B = len(event_B) / len(outcomes)
    print(f"事件B(>4): {event_B}, P(B) = {prob_B}")

    # 事件A∩B：既是偶数又大于4
    event_A_intersect_B = [outcome for outcome in event_A if outcome in event_B]
    prob_A_intersect_B = len(event_A_intersect_B) / len(outcomes)
    print(f"A∩B: {event_A_intersect_B}, P(A∩B) = {prob_A_intersect_B}")

    # 验证加法公式
    prob_A_union_B = prob_A + prob_B - prob_A_intersect_B
    print(f"P(A∪B) = P(A) + P(B) - P(A∩B) = {prob_A_union_B}")

    # 直接计算P(A∪B)
    event_A_union_B = list(set(event_A + event_B))
    prob_A_union_B_direct = len(event_A_union_B) / len(outcomes)
    print(f"直接计算 P(A∪B) = {prob_A_union_B_direct}")

def simulate_coin_flips():
    """模拟抛硬币实验"""
    print("\n=== 抛硬币实验模拟 ===")

    n_flips = 10000
    coin_flips = np.random.choice(['H', 'T'], size=n_flips, p=[0.5, 0.5])

    # 统计正面和反面次数
    counts = Counter(coin_flips)
    print(f"抛硬币 {n_flips} 次:")
    print(f"正面(H): {counts['H']} 次, 概率 ≈ {counts['H']/n_flips:.4f}")
    print(f"反面(T): {counts['T']} 次, 概率 ≈ {counts['T']/n_flips:.4f}")

    # 计算连续正面的最大长度
    max_consecutive_heads = 0
    current_consecutive = 0

    for flip in coin_flips:
        if flip == 'H':
            current_consecutive += 1
            max_consecutive_heads = max(max_consecutive_heads, current_consecutive)
        else:
            current_consecutive = 0

    print(f"最长连续正面: {max_consecutive_heads} 次")

    # 可视化结果
    plt.figure(figsize=(12, 4))

    plt.subplot(121)
    plt.bar(['正面', '反面'], [counts['H'], counts['T']], color=['blue', 'red'])
    plt.title(f'{n_flips}次抛硬币结果')
    plt.ylabel('次数')

    plt.subplot(122)
    # 显示前100次的结果序列
    first_100 = coin_flips[:100]
    heads_positions = [i for i, flip in enumerate(first_100) if flip == 'H']
    tails_positions = [i for i, flip in enumerate(first_100) if flip == 'T']

    plt.scatter(heads_positions, [1]*len(heads_positions),
               c='blue', marker='o', s=20, alpha=0.6, label='正面')
    plt.scatter(tails_positions, [0]*len(tails_positions),
               c='red', marker='x', s=20, alpha=0.6, label='反面')
    plt.yticks([0, 1], ['反面', '正面'])
    plt.xlabel('试验次数')
    plt.title('前100次抛硬币序列')
    plt.legend()
    plt.grid(True, alpha=0.3)

    plt.tight_layout()
    plt.show()

# ==================== 条件概率与贝叶斯定理 ====================

def demonstrate_conditional_probability():
    """演示条件概率"""
    print("\n=== 条件概率 ===")

    # 示例：医疗诊断
    # 假设某种疾病在人群中的发病率为1%
    # 检测试剂的准确率为99%
    # 如果某人检测结果为阳性，他真正患病的概率是多少？

    P_disease = 0.01  # 患病概率
    P_no_disease = 0.99  # 不患病概率
    P_positive_given_disease = 0.99  # 患病情况下检测为阳性的概率
    P_positive_given_no_disease = 0.01  # 不患病情况下检测为阳性的概率（假阳性）

    # 使用全概率公式计算检测为阳性的总概率
    P_positive = (P_positive_given_disease * P_disease +
                  P_positive_given_no_disease * P_no_disease)

    # 使用贝叶斯定理计算检测为阳性时真正患病的概率
    P_disease_given_positive = (P_positive_given_disease * P_disease) / P_positive

    print("医疗诊断示例：")
    print(f"人群患病率: {P_disease:.1%}")
    print(f"检测准确率: {P_positive_given_disease:.1%}")
    print(f"检测为阳性时真正患病的概率: {P_disease_given_positive:.1%}")

    # 直观解释
    print("\n直观解释（假设10000人）：")
    population = 10000
    diseased = population * P_disease
    healthy = population * P_no_disease
    true_positives = diseased * P_positive_given_disease
    false_positives = healthy * P_positive_given_no_disease
    total_positives = true_positives + false_positives

    print(f"总人数: {population}")
    print(f"患病人数: {diseased:.0f}")
    print(f"健康人数: {healthy:.0f}")
    print(f"真阳性: {true_positives:.0f}")
    print(f"假阳性: {false_positives:.0f}")
    print(f"总阳性: {total_positives:.0f}")
    print(f"阳性中真正患病的比例: {true_positives/total_positives:.1%}")

def simulate_conditional_probability():
    """模拟条件概率"""
    print("\n=== 条件概率模拟 ===")

    # 模拟抛掷两个骰子
    n_trials = 100000
    die1 = np.random.randint(1, 7, size=n_trials)
    die2 = np.random.randint(1, 7, size=n_trials)
    sum_dice = die1 + die2

    # 计算P(总和>=10 | 第一个骰子=6)
    event_A = (die1 == 6)  # 第一个骰子为6
    event_B = (sum_dice >= 10)  # 总和>=10
    event_A_and_B = event_A & event_B

    P_A = np.mean(event_A)
    P_B = np.mean(event_B)
    P_A_and_B = np.mean(event_A_and_B)
    P_B_given_A = P_A_and_B / P_A if P_A > 0 else 0

    print(f"P(第一个骰子=6) = {P_A:.4f}")
    print(f"P(总和>=10) = {P_B:.4f}")
    print(f"P(第一个骰子=6 且 总和>=10) = {P_A_and_B:.4f}")
    print(f"P(总和>=10 | 第一个骰子=6) = {P_B_given_A:.4f}")

    # 理论值
    # 如果第一个骰子是6，第二个骰子可以是4,5,6才能使总和>=10
    theoretical_P_B_given_A = 3/6
    print(f"理论值: P(总和>=10 | 第一个骰子=6) = {theoretical_P_B_given_A:.4f}")

# ==================== 随机变量与概率分布 ====================

def demonstrate_discrete_distributions():
    """演示离散概率分布"""
    print("\n=== 离散概率分布 ===")

    # 伯努利分布
    print("1. 伯努利分布（单次硬币抛掷）:")
    p = 0.6
    bernoulli = stats.bernoulli(p)
    print(f"P(X=1) = {p}, P(X=0) = {1-p}")
    print(f"期望: E[X] = {bernoulli.mean():.3f}")
    print(f"方差: Var[X] = {bernoulli.var():.3f}")

    # 二项分布
    print("\n2. 二项分布（10次硬币抛掷，正面次数）:")
    n, p = 10, 0.6
    binomial = stats.binom(n, p)
    print(f"参数: n={n}, p={p}")
    print(f"期望: E[X] = {binomial.mean():.3f}")
    print(f"方差: Var[X] = {binomial.var():.3f}")
    print(f"P(X=5) = {binomial.pmf(5):.4f}")
    print(f"P(X≤5) = {binomial.cdf(5):.4f}")

    # 泊松分布
    print("\n3. 泊松分布（单位时间内事件发生次数）:")
    lam = 3
    poisson = stats.poisson(lam)
    print(f"参数: λ={lam}")
    print(f"期望: E[X] = {poisson.mean():.3f}")
    print(f"方差: Var[X] = {poisson.var():.3f}")
    print(f"P(X=2) = {poisson.pmf(2):.4f}")
    print(f"P(X≤2) = {poisson.cdf(2):.4f}")

    visualize_discrete_distributions()

def visualize_discrete_distributions():
    """可视化离散分布"""
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))

    # 伯努利分布
    ax1 = axes[0, 0]
    p = 0.6
    x = [0, 1]
    y = [1-p, p]
    ax1.bar(x, y, alpha=0.7, color=['blue', 'red'])
    ax1.set_xticks(x)
    ax1.set_xlabel('X')
    ax1.set_ylabel('P(X)')
    ax1.set_title(f'伯努利分布 (p={p})')

    # 二项分布
    ax2 = axes[0, 1]
    n, p = 10, 0.6
    x = np.arange(0, n+1)
    y = stats.binom.pmf(x, n, p)
    ax2.bar(x, y, alpha=0.7)
    ax2.set_xlabel('X')
    ax2.set_ylabel('P(X)')
    ax2.set_title(f'二项分布 (n={n}, p={p})')

    # 泊松分布
    ax3 = axes[1, 0]
    lam = 3
    x = np.arange(0, 15)
    y = stats.poisson.pmf(x, lam)
    ax3.bar(x, y, alpha=0.7)
    ax3.set_xlabel('X')
    ax3.set_ylabel('P(X)')
    ax3.set_title(f'泊松分布 (λ={lam})')

    # 几何分布
    ax4 = axes[1, 1]
    p = 0.3
    x = np.arange(1, 15)
    y = stats.geom.pmf(x, p)
    ax4.bar(x, y, alpha=0.7)
    ax4.set_xlabel('X')
    ax4.set_ylabel('P(X)')
    ax4.set_title(f'几何分布 (p={p})')

    plt.tight_layout()
    plt.show()

def demonstrate_continuous_distributions():
    """演示连续概率分布"""
    print("\n=== 连续概率分布 ===")

    # 正态分布
    print("1. 正态分布:")
    mu, sigma = 0, 1
    normal = stats.norm(mu, sigma)
    print(f"参数: μ={mu}, σ={sigma}")
    print(f"期望: E[X] = {normal.mean():.3f}")
    print(f"方差: Var[X] = {normal.var():.3f}")
    print(f"P(-1≤X≤1) = {normal.cdf(1) - normal.cdf(-1):.4f}")
    print(f"P(|X|≤1.96) = {normal.cdf(1.96) - normal.cdf(-1.96):.4f}")

    # 指数分布
    print("\n2. 指数分布:")
    scale = 2  # 1/λ
    exponential = stats.expon(scale=scale)
    lam = 1/scale
    print(f"参数: λ={lam:.3f}")
    print(f"期望: E[X] = {exponential.mean():.3f}")
    print(f"方差: Var[X] = {exponential.var():.3f}")
    print(f"P(X≤2) = {exponential.cdf(2):.4f}")

    # 均匀分布
    print("\n3. 均匀分布:")
    a, b = 0, 10
    uniform = stats.uniform(a, b-a)
    print(f"参数: a={a}, b={b}")
    print(f"期望: E[X] = {uniform.mean():.3f}")
    print(f"方差: Var[X] = {uniform.var():.3f}")

    visualize_continuous_distributions()

def visualize_continuous_distributions():
    """可视化连续分布"""
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))

    # 正态分布
    ax1 = axes[0, 0]
    x = np.linspace(-4, 4, 1000)
    for mu, sigma, color, label in [(0, 1, 'blue', 'N(0,1)'),
                                    (0, 2, 'red', 'N(0,4)'),
                                    (2, 1, 'green', 'N(2,1)')]:
        y = stats.norm.pdf(x, mu, sigma)
        ax1.plot(x, y, color=color, label=label)
    ax1.set_xlabel('x')
    ax1.set_ylabel('f(x)')
    ax1.set_title('正态分布')
    ax1.legend()

    # 指数分布
    ax2 = axes[0, 1]
    x = np.linspace(0, 10, 1000)
    for lam, color, label in [(0.5, 'blue', 'λ=0.5'),
                             (1, 'red', 'λ=1'),
                             (2, 'green', 'λ=2')]:
        y = stats.expon.pdf(x, scale=1/lam)
        ax2.plot(x, y, color=color, label=label)
    ax2.set_xlabel('x')
    ax2.set_ylabel('f(x)')
    ax2.set_title('指数分布')
    ax2.legend()

    # 均匀分布
    ax3 = axes[1, 0]
    x = np.linspace(-1, 11, 1000)
    y = stats.uniform.pdf(x, 0, 10)
    ax3.plot(x, y, 'blue', linewidth=2)
    ax3.fill_between(x, 0, y, alpha=0.3)
    ax3.set_xlabel('x')
    ax3.set_ylabel('f(x)')
    ax3.set_title('均匀分布 U(0,10)')

    # 伽马分布
    ax4 = axes[1, 1]
    x = np.linspace(0, 15, 1000)
    for alpha, beta, color, label in [(1, 1, 'blue', 'α=1,β=1'),
                                      (2, 1, 'red', 'α=2,β=1'),
                                      (3, 1, 'green', 'α=3,β=1')]:
        y = stats.gamma.pdf(x, alpha, scale=1/beta)
        ax4.plot(x, y, color=color, label=label)
    ax4.set_xlabel('x')
    ax4.set_ylabel('f(x)')
    ax4.set_title('伽马分布')
    ax4.legend()

    plt.tight_layout()
    plt.show()

# ==================== 大数定律与中心极限定理 ====================

def demonstrate_law_of_large_numbers():
    """演示大数定律"""
    print("\n=== 大数定律演示 ===")

    n_trials = 10000
    # 生成一系列随机数（均匀分布在[0,1]）
    random_numbers = np.random.random(n_trials)

    # 计算累积均值
    cumulative_means = np.cumsum(random_numbers) / np.arange(1, n_trials + 1)

    # 理论期望值
    theoretical_mean = 0.5

    print(f"理论期望值: {theoretical_mean}")
    print(f"最后1000次试验的平均值: {np.mean(random_numbers[-1000:]):.4f}")
    print(f"所有试验的平均值: {np.mean(random_numbers):.4f}")

    # 可视化
    plt.figure(figsize=(10, 6))
    plt.plot(cumulative_means, 'b-', alpha=0.7, label='累积均值')
    plt.axhline(y=theoretical_mean, color='r', linestyle='--',
                label=f'理论期望值 ({theoretical_mean})')

    # 添加收敛区间
    plt.fill_between(range(n_trials),
                    theoretical_mean - 0.01,
                    theoretical_mean + 0.01,
                    alpha=0.2, color='green',
                    label='±0.01 容差区间')

    plt.xlabel('试验次数')
    plt.ylabel('样本均值')
    plt.title('大数定律：样本均值收敛于期望值')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.ylim(0.45, 0.55)
    plt.show()

def demonstrate_central_limit_theorem():
    """演示中心极限定理"""
    print("\n=== 中心极限定理演示 ===")

    # 从指数分布中抽样，观察样本均值的分布
    sample_sizes = [1, 5, 30, 100]
    n_samples = 1000
    true_lambda = 2  # 指数分布参数
    true_mean = 1/true_lambda
    true_var = 1/(true_lambda**2)

    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    axes = axes.flatten()

    for idx, sample_size in enumerate(sample_sizes):
        # 生成样本
        sample_means = []
        for _ in range(n_samples):
            sample = np.random.exponential(scale=1/true_lambda, size=sample_size)
            sample_means.append(np.mean(sample))

        sample_means = np.array(sample_means)

        # 计算理论正态分布参数
        theoretical_mean = true_mean
        theoretical_std = np.sqrt(true_var / sample_size)

        # 绘制直方图
        ax = axes[idx]
        ax.hist(sample_means, bins=30, density=True, alpha=0.7,
                color='skyblue', label='样本均值分布')

        # 绘制理论正态分布
        x = np.linspace(sample_means.min(), sample_means.max(), 1000)
        y = stats.norm.pdf(x, theoretical_mean, theoretical_std)
        ax.plot(x, y, 'r-', linewidth=2, label='理论正态分布')

        ax.set_xlabel('样本均值')
        ax.set_ylabel('密度')
        ax.set_title(f'样本大小 = {sample_size}')
        ax.legend()

        # 添加统计信息
        ax.text(0.05, 0.95, f'样本均值: {np.mean(sample_means):.3f}\n'
                           f'理论均值: {theoretical_mean:.3f}\n'
                           f'样本标准差: {np.std(sample_means):.3f}\n'
                           f'理论标准差: {theoretical_std:.3f}',
                transform=ax.transAxes, verticalalignment='top',
                bbox=dict(boxstyle='round', facecolor='white', alpha=0.8))

    plt.suptitle('中心极限定理：样本均值的分布趋近于正态分布')
    plt.tight_layout()
    plt.show()

# ==================== 统计推断基础 ====================

def demonstrate_hypothesis_testing():
    """演示假设检验"""
    print("\n=== 假设检验演示 ===")

    # 示例：检验新药是否有效
    # H0: μ = 100 (药物无效)
    # H1: μ > 100 (药物有效)

    # 生成模拟数据
    np.random.seed(42)
    control_group = np.random.normal(100, 15, 50)  # 对照组
    treatment_group = np.random.normal(105, 15, 50)  # 治疗组

    control_mean = np.mean(control_group)
    treatment_mean = np.mean(treatment_group)
    control_std = np.std(control_group, ddof=1)
    treatment_std = np.std(treatment_group, ddof=1)

    print(f"对照组: 均值 = {control_mean:.2f}, 标准差 = {control_std:.2f}")
    print(f"治疗组: 均值 = {treatment_mean:.2f}, 标准差 = {treatment_std:.2f}")

    # 两样本t检验
    t_stat, p_value = stats.ttest_ind(treatment_group, control_group,
                                     alternative='greater')

    print(f"\nt统计量: {t_stat:.4f}")
    print(f"p值: {p_value:.4f}")

    alpha = 0.05
    if p_value < alpha:
        print(f"p值 < {alpha}, 拒绝原假设H0")
        print("结论：药物具有显著效果")
    else:
        print(f"p值 ≥ {alpha}, 不能拒绝原假设H0")
        print("结论：没有足够证据表明药物有效")

    visualize_hypothesis_testing(control_group, treatment_group)

def visualize_hypothesis_testing(control_group, treatment_group):
    """可视化假设检验结果"""
    plt.figure(figsize=(12, 5))

    # 箱线图比较
    plt.subplot(121)
    data = [control_group, treatment_group]
    labels = ['对照组', '治疗组']
    colors = ['lightblue', 'lightcoral']

    bp = plt.boxplot(data, labels=labels, patch_artist=True)
    for patch, color in zip(bp['boxes'], colors):
        patch.set_facecolor(color)

    plt.ylabel('测量值')
    plt.title('两组数据分布比较')

    # 添加均值线
    means = [np.mean(control_group), np.mean(treatment_group)]
    plt.plot([1, 2], means, 'ro-', linewidth=2, markersize=8, label='均值')
    plt.legend()

    # 直方图比较
    plt.subplot(122)
    plt.hist(control_group, bins=15, alpha=0.7, color='blue',
             label='对照组', density=True)
    plt.hist(treatment_group, bins=15, alpha=0.7, color='red',
             label='治疗组', density=True)

    plt.xlabel('测量值')
    plt.ylabel('密度')
    plt.title('数据分布密度')
    plt.legend()

    plt.tight_layout()
    plt.show()

# ==================== 蒙特卡罗模拟 ====================

def monte_carlo_pi():
    """用蒙特卡罗方法计算π"""
    print("\n=== 蒙特卡罗方法计算π ===")

    n_points = 10000
    x = np.random.random(n_points) * 2 - 1  # [-1, 1]
    y = np.random.random(n_points) * 2 - 1  # [-1, 1]

    # 判断点是否在单位圆内
    distances = np.sqrt(x**2 + y**2)
    inside_circle = distances <= 1

    n_inside = np.sum(inside_circle)
    pi_estimate = 4 * n_inside / n_points

    print(f"总点数: {n_points}")
    print(f"圆内点数: {n_inside}")
    print(f"π的估计值: {pi_estimate:.6f}")
    print(f"真实π值: {np.pi:.6f}")
    print(f"相对误差: {abs(pi_estimate - np.pi)/np.pi:.4%}")

    # 可视化
    plt.figure(figsize=(8, 8))
    plt.scatter(x[inside_circle], y[inside_circle], c='blue', s=1, alpha=0.5, label='圆内')
    plt.scatter(x[~inside_circle], y[~inside_circle], c='red', s=1, alpha=0.5, label='圆外')

    # 绘制单位圆
    circle = plt.Circle((0, 0), 1, fill=False, color='black', linewidth=2)
    plt.gca().add_patch(circle)

    # 绘制正方形边界
    plt.plot([-1, 1, 1, -1, -1], [-1, -1, 1, 1, -1], 'k-', linewidth=2)

    plt.xlim(-1.2, 1.2)
    plt.ylim(-1.2, 1.2)
    plt.axis('equal')
    plt.legend()
    plt.title(f'蒙特卡罗方法估计π = {pi_estimate:.6f}')
    plt.show()

def monte_carlo_integration():
    """蒙特卡罗积分示例"""
    print("\n=== 蒙特卡罗积分 ===")

    # 计算函数 f(x) = x² 在 [0, 1] 上的积分
    # 真实值是 1/3

    def f(x):
        return x**2

    n_samples = 100000
    x_random = np.random.random(n_samples)
    y_random = np.random.random(n_samples)

    # 计算函数在随机点的值
    y_function = f(x_random)

    # 计算在曲线下方的点数
    under_curve = y_random <= y_function
    n_under = np.sum(under_curve)

    # 估计积分值
    integral_estimate = n_under / n_samples

    print(f"积分 ∫[0,1] x² dx 的估计值: {integral_estimate:.6f}")
    print(f"真实值: {1/3:.6f}")
    print(f"相对误差: {abs(integral_estimate - 1/3)/(1/3):.4%}")

    # 可视化
    plt.figure(figsize=(10, 6))

    x = np.linspace(0, 1, 1000)
    y = f(x)

    plt.plot(x, y, 'b-', linewidth=2, label='f(x) = x²')

    # 显示随机点
    plt.scatter(x_random[under_curve], y_random[under_curve],
               c='green', s=1, alpha=0.3, label='曲线下')
    plt.scatter(x_random[~under_curve], y_random[~under_curve],
               c='red', s=1, alpha=0.3, label='曲线上')

    # 填充积分区域
    plt.fill_between(x, 0, y, alpha=0.2, color='blue')

    plt.xlim(0, 1)
    plt.ylim(0, 1.2)
    plt.xlabel('x')
    plt.ylabel('y')
    plt.legend()
    plt.title(f'蒙特卡罗积分: ∫[0,1] x² dx ≈ {integral_estimate:.6f}')
    plt.grid(True, alpha=0.3)
    plt.show()

# ==================== 主函数 ====================

def main():
    """主函数，运行所有演示"""
    print("概率论学习演示程序")
    print("=" * 50)

    # 运行各个模块
    demonstrate_basic_probability()
    simulate_coin_flips()
    demonstrate_conditional_probability()
    simulate_conditional_probability()
    demonstrate_discrete_distributions()
    demonstrate_continuous_distributions()
    demonstrate_law_of_large_numbers()
    demonstrate_central_limit_theorem()
    demonstrate_hypothesis_testing()
    monte_carlo_pi()
    monte_carlo_integration()

    print("\n演示完成！")
    print("建议：")
    print("1. 尝试改变参数，观察分布和结果的变化")
    print("2. 理解各种概率公式的含义和应用场景")
    print("3. 通过模拟加深对概率概念的理解")

if __name__ == "__main__":
    main()