#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
考研英语一学习系统 - 主程序
Main entry point for the Graduate English Exam Learning System
"""

import sys
import os
from pathlib import Path

# 添加当前目录到路径，以便导入子模块
sys.path.insert(0, str(Path(__file__).parent))


def print_banner():
    """打印欢迎横幅"""
    banner = """
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║          考研英语一学习系统                               ║
║          Graduate English Exam Learning System            ║
║                                                           ║
║          版本：v1.0                                       ║
║          专注：词汇方法论 + 阅读策略                      ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
"""
    print(banner)


def print_main_menu():
    """打印主菜单"""
    menu = """
┌─────────────────────────────────────┐
│         主菜单 Main Menu            │
├─────────────────────────────────────┤
│  1. 词汇学习模式                     │
│     Vocabulary Learning Mode        │
│                                     │
│  2. 阅读策略查看                     │
│     Reading Strategies Overview     │
│                                     │
│  3. 学习统计                         │
│     Study Statistics                │
│                                     │
│  4. 帮助文档                         │
│     Help Documentation              │
│                                     │
│  0. 退出系统                         │
│     Exit                            │
└─────────────────────────────────────┘
"""
    print(menu)


def vocabulary_mode():
    """词汇学习模式"""
    try:
        from vocabulary.vocabulary_tool import VocabularyTool
        
        print("\n" + "="*60)
        print("进入词汇学习模式...")
        print("="*60 + "\n")
        
        tool = VocabularyTool()
        tool.run()
        
    except ImportError as e:
        print(f"\n[错误] 无法导入词汇工具模块")
        print(f"   详细信息：{e}")
        print("   请确保 vocabulary/vocabulary_tool.py 文件存在")
    except Exception as e:
        print(f"\n[错误] 发生错误：{e}")


def reading_overview():
    """阅读策略概览"""
    print("\n" + "="*60)
    print("阅读理解策略概览")
    print("="*60 + "\n")
    
    overview = """
考研英语一阅读理解核心策略

一、阅读三步法
   1. Skimming（略读）- 3-5分钟把握文章主旨
   2. Scanning（扫读）- 根据题目快速定位关键信息
   3. Close Reading（精读）- 仔细理解定位段落

二、六大题型及策略
   1. 主旨大意题 - 关注首尾段，识别主题句
   2. 细节题 - 关键词定位+同义替换识别
   3. 推理判断题 - 合理推理，避免过度推断
   4. 词义猜测题 - 利用上下文和词根词缀
   5. 态度观点题 - 识别情感词汇（positive/negative/critical等）
   6. 结构题（例证题）- 理解例子服务于论点

三、长难句分析步骤
   1. 找主干（主谓宾）
   2. 识别从句（定语从句、状语从句等）
   3. 剥离修饰成分
   4. 利用标点符号断句

四、时间管理
   - 每篇文章：15-18分钟
   - 阅读文章：5-7分钟
   - 做题：8-10分钟
   - 检查：1-2分钟

五、常见干扰项特征
   [x] 偷换概念（张冠李戴）
   [x] 以偏概全（扩大或缩小范围）
   [x] 无中生有（原文未提及）
   [x] 过度推断（推理过远）

六、真题练习四步法
   第一步：限时完成（模拟考试）
   第二步：对答案，标注错题
   第三步：逐句翻译全文
   第四步：总结错因，整理笔记

详细内容请查看：
   - reading/reading_strategies_guide.md（策略详解）
   - reading/question_types_analysis.md（题型分析）
   - reading/article_structure_guide.md（文章结构）
"""
    print(overview)
    input("\n按 Enter 键返回主菜单...")


def show_statistics():
    """显示学习统计"""
    print("\n" + "="*60)
    print("学习统计")
    print("="*60 + "\n")
    
    try:
        from vocabulary.vocabulary_tool import VocabularyTool
        tool = VocabularyTool()
        tool.show_statistics()
    except ImportError:
        print("[警告] 词汇工具模块未找到，无法显示统计数据")
    except Exception as e:
        print(f"[错误] {e}")
    
    input("\n按 Enter 键返回主菜单...")


def show_help():
    """显示帮助文档"""
    help_text = """
+-----------------------------------------------------------+
|                    帮助文档                               |
+-----------------------------------------------------------+

【命令行参数】

  python main.py              # 交互式菜单模式（推荐）
  python main.py --vocab      # 直接进入词汇学习
  python main.py --reading    # 查看阅读策略
  python main.py --help       # 显示帮助信息

【词汇学习功能】

  1. 随机测试：从词库中随机抽取单词进行测试
  2. 艾宾浩斯复习：根据记忆曲线自动提醒复习
  3. 错题本：自动记录错误单词，强化练习
  4. 学习统计：查看学习进度和正确率
  5. 词汇搜索：快速查找单词释义

【数据文件位置】

  - 词汇数据：user_data/words.json
  - 错题本：user_data/wrong_words.json
  - 学习记录：user_data/learning_log.json

【文档资源】

  vocabulary/
     - vocabulary_guide.md      词汇学习方法论
     - word_groups.md           词根词缀分类
     - high_frequency_words.md  高频词汇主题分类

  reading/
     - reading_strategies_guide.md     阅读策略指南
     - question_types_analysis.md      题型解析
     - article_structure_guide.md      文章结构识别

  docs/
     - exam_analysis.md          考试分析与备考建议

【学习建议】

  * 每日词汇学习：20-30个新词 + 复习旧词
  * 遵循艾宾浩斯曲线：1天、2天、4天、7天、15天复习
  * 阅读训练：前期每天1篇精读，中期2篇限时，后期4篇模拟
  * 真题为王：近10年真题至少做3遍

【技术支持】

  遇到问题？
  1. 检查 Python 版本是否为 3.7+
  2. 确保所有文件完整（尤其是 vocabulary_tool.py）
  3. 查看 user_data/ 目录权限

"""
    print(help_text)
    input("\n按 Enter 键返回主菜单...")


def main_menu():
    """主菜单循环"""
    while True:
        print_main_menu()
        
        try:
            choice = input("请选择功能（输入数字）：").strip()
            
            if choice == '1':
                vocabulary_mode()
            elif choice == '2':
                reading_overview()
            elif choice == '3':
                show_statistics()
            elif choice == '4':
                show_help()
            elif choice == '0':
                print("\n感谢使用！祝考研顺利！")
                print("Goodbye and good luck with your exam!\n")
                sys.exit(0)
            else:
                print("\n[警告] 无效选择，请输入 0-4 之间的数字")
                input("按 Enter 键继续...")
        
        except KeyboardInterrupt:
            print("\n\n检测到中断信号，正在退出...")
            sys.exit(0)
        except EOFError:
            print("\n\n检测到输入结束，正在退出...")
            sys.exit(0)


def main():
    """主函数"""
    # 检查命令行参数
    if len(sys.argv) > 1:
        arg = sys.argv[1].lower()
        
        if arg in ['--help', '-h', 'help']:
            print_banner()
            show_help()
            return
        
        elif arg in ['--vocab', '-v', 'vocab']:
            print_banner()
            vocabulary_mode()
            return
        
        elif arg in ['--reading', '-r', 'reading']:
            print_banner()
            reading_overview()
            return
        
        else:
            print(f"[错误] 未知参数：{arg}")
            print("使用 --help 查看帮助信息")
            sys.exit(1)
    
    # 交互式菜单模式
    print_banner()
    print("\n欢迎使用考研英语一学习系统！")
    print("本系统提供科学的词汇学习方法和系统的阅读策略指导。\n")
    input("按 Enter 键进入主菜单...")
    
    main_menu()


if __name__ == "__main__":
    main()
