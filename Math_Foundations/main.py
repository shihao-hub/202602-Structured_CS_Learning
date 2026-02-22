"""
数学基础学习项目 - 主程序入口

本项目包含三个核心数学模块：
1. 高等数学 (Calculus)
2. 线性代数 (Linear Algebra)
3. 概率论与数理统计 (Probability & Statistics)

运行此文件可以选择性地运行各个模块的示例代码。
"""

import sys
from pathlib import Path


def print_banner():
    """打印欢迎横幅"""
    print("=" * 60)
    print("欢迎来到数学基础学习项目!")
    print("=" * 60)
    print()
    print("本项目包含以下模块:")
    print("1. 高等数学 (Calculus)")
    print("2. 线性代数 (Linear Algebra)")
    print("3. 概率论与数理统计 (Probability & Statistics)")
    print()


def print_menu():
    """打印选择菜单"""
    print("请选择要运行的模块:")
    print("1 - 高等数学示例")
    print("2 - 线性代数示例")
    print("3 - 概率论示例")
    print("4 - 运行所有示例")
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


def run_all_examples():
    """运行所有模块的示例"""
    print("\n正在运行所有模块的示例...")
    print("注意: 可视化图表会依次显示，请关闭当前图表后继续。")
    print("-" * 60)
    
    run_calculus_examples()
    run_linear_algebra_examples()
    run_probability_examples()
    
    print("\n" + "=" * 60)
    print("所有示例运行完成!")
    print("=" * 60)


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
    print("5. 尝试解决课后练习题")
    print()
    print("推荐学习顺序:")
    print("  高等数学 → 线性代数 → 概率论")
    print()


def interactive_mode():
    """交互式运行模式"""
    while True:
        print_menu()
        choice = input("请输入选项 (0-4): ").strip()
        
        if choice == '0':
            print("\n感谢使用，再见! Happy Learning! 🚀")
            break
        elif choice == '1':
            run_calculus_examples()
        elif choice == '2':
            run_linear_algebra_examples()
        elif choice == '3':
            run_probability_examples()
        elif choice == '4':
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
        elif arg in ['--all', '-a', '4']:
            run_all_examples()
        elif arg in ['--help', '-h']:
            print("使用方法:")
            print("  python main.py              # 交互模式")
            print("  python main.py --calculus   # 运行高等数学示例")
            print("  python main.py --linear-algebra  # 运行线性代数示例")
            print("  python main.py --probability     # 运行概率论示例")
            print("  python main.py --all        # 运行所有示例")
        else:
            print(f"未知参数: {arg}")
            print("使用 --help 查看帮助信息")
    else:
        # 交互模式
        print_learning_tips()
        interactive_mode()


if __name__ == "__main__":
    main()
