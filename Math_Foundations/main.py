"""
数学基础学习项目 - 主程序入口

本项目包含六个核心数学模块：
基础模块：
1. 高等数学 (Calculus)
2. 线性代数 (Linear Algebra)
3. 概率论与数理统计 (Probability & Statistics)

进阶模块（考研数学一）：
4. 高等数学进阶 (Calculus Advanced)
5. 线性代数进阶 (Linear Algebra Advanced)
6. 概率统计进阶 (Probability Advanced)

运行此文件可以选择性地运行各个模块的示例代码。
"""

import sys
from pathlib import Path


def print_banner():
    """打印欢迎横幅"""
    print("=" * 60)
    print("欢迎来到数学基础学习项目!")
    print("考研数学一 · 基础 + 进阶全覆盖")
    print("=" * 60)
    print()
    print("本项目包含以下模块:")
    print()
    print("【基础模块】")
    print("1. 高等数学 (Calculus)")
    print("2. 线性代数 (Linear Algebra)")
    print("3. 概率论与数理统计 (Probability & Statistics)")
    print()
    print("【进阶模块 - 考研数学一】")
    print("4. 高等数学进阶 (曲线曲面积分/场论/傅里叶)")
    print("5. 线性代数进阶 (二次型/向量空间/谱分解)")
    print("6. 概率统计进阶 (抽样分布/参数估计/假设检验)")
    print()


def print_menu():
    """打印选择菜单"""
    print("请选择要运行的模块:")
    print("1 - 高等数学示例")
    print("2 - 线性代数示例")
    print("3 - 概率论示例")
    print("4 - 高等数学进阶示例")
    print("5 - 线性代数进阶示例")
    print("6 - 概率统计进阶示例")
    print("7 - 运行所有基础模块")
    print("8 - 运行所有进阶模块")
    print("9 - 运行全部示例")
    print("0 - 退出")
    print()


def run_calculus_examples():
    """运行高等数学示例"""
    print("\n" + "=" * 60)
    print("【模块 1: 高等数学】")
    print("=" * 60)
    try:
        from calculus.calculus_examples import main as calculus_main
        calculus_main()
    except ImportError as e:
        print(f"导入高等数学模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib sympy scipy")
    except Exception as e:
        print(f"运行高等数学示例时出错: {e}")


def run_linear_algebra_examples():
    """运行线性代数示例"""
    print("\n" + "=" * 60)
    print("【模块 2: 线性代数】")
    print("=" * 60)
    try:
        from linear_algebra.linear_algebra_examples import main as la_main
        la_main()
    except ImportError as e:
        print(f"导入线性代数模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib scipy seaborn")
    except Exception as e:
        print(f"运行线性代数示例时出错: {e}")


def run_probability_examples():
    """运行概率论示例"""
    print("\n" + "=" * 60)
    print("【模块 3: 概率论与数理统计】")
    print("=" * 60)
    try:
        from probability.probability_examples import main as prob_main
        prob_main()
    except ImportError as e:
        print(f"导入概率论模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib scipy seaborn")
    except Exception as e:
        print(f"运行概率论示例时出错: {e}")


def run_calculus_advanced_examples():
    """运行高等数学进阶示例"""
    print("\n" + "=" * 60)
    print("【模块 4: 高等数学进阶】")
    print("=" * 60)
    try:
        from calculus_advanced.calculus_advanced_examples import main as ca_main
        ca_main()
    except ImportError as e:
        print(f"导入高等数学进阶模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib sympy scipy")
    except Exception as e:
        print(f"运行高等数学进阶示例时出错: {e}")


def run_linear_algebra_advanced_examples():
    """运行线性代数进阶示例"""
    print("\n" + "=" * 60)
    print("【模块 5: 线性代数进阶】")
    print("=" * 60)
    try:
        from linear_algebra_advanced.linear_algebra_advanced_examples import main as laa_main
        laa_main()
    except ImportError as e:
        print(f"导入线性代数进阶模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib scipy")
    except Exception as e:
        print(f"运行线性代数进阶示例时出错: {e}")


def run_probability_advanced_examples():
    """运行概率统计进阶示例"""
    print("\n" + "=" * 60)
    print("【模块 6: 概率统计进阶】")
    print("=" * 60)
    try:
        from probability_advanced.probability_advanced_examples import main as pa_main
        pa_main()
    except ImportError as e:
        print(f"导入概率统计进阶模块失败: {e}")
        print("请确保已安装所需依赖: pip install numpy matplotlib scipy")
    except Exception as e:
        print(f"运行概率统计进阶示例时出错: {e}")


def run_all_basic():
    """运行所有基础模块"""
    print("\n正在运行所有基础模块的示例...")
    print("注意: 可视化图表会依次显示，请关闭当前图表后继续。")
    print("-" * 60)
    run_calculus_examples()
    run_linear_algebra_examples()
    run_probability_examples()
    print("\n" + "=" * 60)
    print("所有基础模块运行完成!")
    print("=" * 60)


def run_all_advanced():
    """运行所有进阶模块"""
    print("\n正在运行所有进阶模块的示例...")
    print("注意: 可视化图表会依次显示，请关闭当前图表后继续。")
    print("-" * 60)
    run_calculus_advanced_examples()
    run_linear_algebra_advanced_examples()
    run_probability_advanced_examples()
    print("\n" + "=" * 60)
    print("所有进阶模块运行完成!")
    print("=" * 60)


def run_all_examples():
    """运行全部示例"""
    run_all_basic()
    run_all_advanced()


def check_dependencies():
    """检查必要的依赖是否已安装"""
    dependencies = ['numpy', 'matplotlib', 'scipy', 'sympy', 'seaborn']
    missing = []

    for dep in dependencies:
        try:
            __import__(dep)
        except ImportError:
            missing.append(dep)

    if missing:
        print("警告: 以下依赖未安装:")
        for dep in missing:
            print(f"  - {dep}")
        print()
        print("请运行以下命令安装依赖:")
        print(f"  pip install {' '.join(missing)}")
        print()
        return False
    return True


def print_learning_tips():
    """打印学习建议"""
    print("\n学习建议:")
    print("-" * 40)
    print("1. 先阅读各模块的理论文档 (*_guide.md)")
    print("2. 运行示例代码，观察可视化效果")
    print("3. 修改代码参数，探索结果变化")
    print("4. 理解公式背后的几何直观")
    print("5. 结合考研真题练习")
    print()
    print("推荐学习顺序:")
    print("  基础: 高等数学 → 线性代数 → 概率论")
    print("  进阶: 高数进阶 → 线代进阶 → 概率进阶")
    print()


def interactive_mode():
    """交互式运行模式"""
    while True:
        print_menu()
        choice = input("请输入选项 (0-9): ").strip()

        if choice == '0':
            print("\n感谢使用，再见! Happy Learning!")
            break
        elif choice == '1':
            run_calculus_examples()
        elif choice == '2':
            run_linear_algebra_examples()
        elif choice == '3':
            run_probability_examples()
        elif choice == '4':
            run_calculus_advanced_examples()
        elif choice == '5':
            run_linear_algebra_advanced_examples()
        elif choice == '6':
            run_probability_advanced_examples()
        elif choice == '7':
            run_all_basic()
        elif choice == '8':
            run_all_advanced()
        elif choice == '9':
            run_all_examples()
        else:
            print("无效输入，请重新选择。")

        print()


def main():
    """主函数"""
    # 添加当前目录到路径，以便导入子模块
    current_dir = Path(__file__).parent
    if str(current_dir) not in sys.path:
        sys.path.insert(0, str(current_dir))

    print_banner()

    # 检查依赖
    if not check_dependencies():
        return

    # 解析命令行参数
    if len(sys.argv) > 1:
        arg = sys.argv[1].lower()
        if arg in ['--calculus', '-c', '1']:
            run_calculus_examples()
        elif arg in ['--linear-algebra', '-l', '2']:
            run_linear_algebra_examples()
        elif arg in ['--probability', '-p', '3']:
            run_probability_examples()
        elif arg in ['--calculus-adv', '-ca', '4']:
            run_calculus_advanced_examples()
        elif arg in ['--la-adv', '-la', '5']:
            run_linear_algebra_advanced_examples()
        elif arg in ['--prob-adv', '-pa', '6']:
            run_probability_advanced_examples()
        elif arg in ['--basic', '-b', '7']:
            run_all_basic()
        elif arg in ['--advanced', '-adv', '8']:
            run_all_advanced()
        elif arg in ['--all', '-a', '9']:
            run_all_examples()
        elif arg in ['--help', '-h']:
            print("使用方法:")
            print("  python main.py                    # 交互模式")
            print()
            print("  基础模块:")
            print("  python main.py --calculus          # 高等数学")
            print("  python main.py --linear-algebra    # 线性代数")
            print("  python main.py --probability       # 概率论")
            print()
            print("  进阶模块 (考研数学一):")
            print("  python main.py --calculus-adv      # 高数进阶")
            print("  python main.py --la-adv            # 线代进阶")
            print("  python main.py --prob-adv          # 概率进阶")
            print()
            print("  批量运行:")
            print("  python main.py --basic             # 所有基础")
            print("  python main.py --advanced           # 所有进阶")
            print("  python main.py --all               # 全部")
        else:
            print(f"未知参数: {arg}")
            print("使用 --help 查看帮助信息")
    else:
        # 交互模式
        print_learning_tips()
        interactive_mode()


if __name__ == "__main__":
    main()
