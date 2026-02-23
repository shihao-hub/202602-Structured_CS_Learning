# 数学基础学习资料

面向考研数学一的系统化数学学习项目，使用Python实现，包含基础和进阶两个层次。

## 文件结构

```
Math_Foundations/
├── README.md                              # 本文件
├── main.py                                # 主程序入口
├── calculus/                              # 高等数学（基础）
│   ├── calculus_guide.md                  # 学习指南
│   └── calculus_examples.py               # 代码示例
├── linear_algebra/                        # 线性代数（基础）
│   ├── linear_algebra_guide.md            # 学习指南
│   └── linear_algebra_examples.py         # 代码示例
├── probability/                           # 概率论（基础）
│   ├── probability_guide.md               # 学习指南
│   └── probability_examples.py            # 代码示例
├── calculus_advanced/                     # 高等数学（进阶·数学一）
│   ├── README.md
│   ├── calculus_advanced_guide.md         # 曲线曲面积分/场论/傅里叶
│   └── calculus_advanced_examples.py      # 代码示例
├── linear_algebra_advanced/               # 线性代数（进阶·数学一）
│   ├── README.md
│   ├── linear_algebra_advanced_guide.md   # 二次型/向量空间/谱分解
│   └── linear_algebra_advanced_examples.py # 代码示例
├── probability_advanced/                  # 概率统计（进阶·数学一）
│   ├── README.md
│   ├── probability_advanced_guide.md      # 抽样分布/参数估计/假设检验
│   └── probability_advanced_examples.py   # 代码示例
└── docs/
    └── learning_summary.md                # 学习总结
```

## 环境要求

```bash
pip install numpy matplotlib scipy sympy seaborn
```

## 使用方法

### 交互模式
```bash
cd Math_Foundations
python main.py
```

### 命令行参数
```bash
# 基础模块
python main.py --calculus          # 高等数学
python main.py --linear-algebra    # 线性代数
python main.py --probability       # 概率论

# 进阶模块（考研数学一）
python main.py --calculus-adv      # 高数进阶
python main.py --la-adv            # 线代进阶
python main.py --prob-adv          # 概率进阶

# 批量运行
python main.py --basic             # 所有基础
python main.py --advanced          # 所有进阶
python main.py --all               # 全部
```

## 学习内容

### 基础模块
| 模块 | 内容 |
|------|------|
| 高等数学 | 极限与连续、导数与微分、积分、级数、多元微积分、常微分方程 |
| 线性代数 | 矩阵运算、线性方程组、向量空间、特征值与特征向量、矩阵分解 |
| 概率论 | 基础概率、条件概率、概率分布、大数定律、假设检验、蒙特卡罗 |

### 进阶模块（考研数学一专用）
| 模块 | 内容 | 考点权重 |
|------|------|---------|
| 高数进阶 | 曲线积分(第一型/第二型)、曲面积分、格林/高斯/斯托克斯公式、场论、傅里叶级数 | ★★★ |
| 线代进阶 | 二次型(标准形/正定)、相似对角化、谱分解、向量空间、四个基本子空间 | ★★★ |
| 概率进阶 | 三大抽样分布(χ²/t/F)、参数估计(MLE/矩估计)、假设检验、一元线性回归 | ★★★ |

## 考研数学一学习路线

```
第一阶段（基础）:
  高等数学基础 → 线性代数基础 → 概率论基础

第二阶段（进阶）:
  高数进阶(曲线曲面积分/场论) → 线代进阶(二次型/向量空间) → 概率进阶(数理统计)

第三阶段（强化）:
  结合真题练习 → 重点突破薄弱环节 → 模拟冲刺
```

## 推荐书籍
- 《高等数学》同济大学
- 《线性代数及其应用》Gilbert Strang
- 《概率论与数理统计》浙大版
