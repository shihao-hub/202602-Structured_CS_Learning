#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
考研英语词汇学习工具
Vocabulary Learning Tool for Graduate English Exam

功能 Features:
1. 随机词汇测试 Random Quiz
2. 艾宾浩斯复习提醒 Ebbinghaus Review Reminder
3. 错题本管理 Wrong Answer Notebook
4. 学习进度统计 Learning Statistics
5. 词汇搜索 Word Search
6. 数据导入/导出 Data Import/Export
"""

import json
import random
import os
from datetime import datetime, timedelta
from pathlib import Path


class VocabularyTool:
    """词汇学习工具类"""
    
    def __init__(self, data_dir="user_data"):
        """初始化词汇工具"""
        self.data_dir = Path(__file__).parent.parent / data_dir
        self.data_dir.mkdir(parents=True, exist_ok=True)
        
        self.words_file = self.data_dir / "words.json"
        self.wrong_file = self.data_dir / "wrong_words.json"
        self.log_file = self.data_dir / "learning_log.json"
        
        self.words = self.load_words()
        self.wrong_words = self.load_wrong_words()
        self.learning_log = self.load_learning_log()
    
    def load_words(self):
        """加载词汇库"""
        if self.words_file.exists():
            with open(self.words_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        else:
            # 如果文件不存在,创建示例词汇库
            sample_words = self.create_sample_words()
            self.save_words(sample_words)
            return sample_words
    
    def create_sample_words(self):
        """创建示例词汇库"""
        return [
            {"english": "abandon", "chinese": "放弃、抛弃", "learned_date": None, "review_count": 0},
            {"english": "abstract", "chinese": "抽象的、摘要", "learned_date": None, "review_count": 0},
            {"english": "achieve", "chinese": "实现、达到", "learned_date": None, "review_count": 0},
            {"english": "adapt", "chinese": "适应、改编", "learned_date": None, "review_count": 0},
            {"english": "advocate", "chinese": "提倡、拥护", "learned_date": None, "review_count": 0},
            {"english": "allocate", "chinese": "分配、配置", "learned_date": None, "review_count": 0},
            {"english": "analyze", "chinese": "分析", "learned_date": None, "review_count": 0},
            {"english": "apparent", "chinese": "明显的", "learned_date": None, "review_count": 0},
            {"english": "approach", "chinese": "方法、接近", "learned_date": None, "review_count": 0},
            {"english": "appropriate", "chinese": "适当的、拨款", "learned_date": None, "review_count": 0},
            {"english": "aspect", "chinese": "方面、外观", "learned_date": None, "review_count": 0},
            {"english": "assess", "chinese": "评估", "learned_date": None, "review_count": 0},
            {"english": "assume", "chinese": "假设、承担", "learned_date": None, "review_count": 0},
            {"english": "attach", "chinese": "附加、依附", "learned_date": None, "review_count": 0},
            {"english": "attitude", "chinese": "态度", "learned_date": None, "review_count": 0},
            {"english": "attribute", "chinese": "归因于、属性", "learned_date": None, "review_count": 0},
            {"english": "authority", "chinese": "权威、当局", "learned_date": None, "review_count": 0},
            {"english": "available", "chinese": "可获得的", "learned_date": None, "review_count": 0},
            {"english": "benefit", "chinese": "利益、有益于", "learned_date": None, "review_count": 0},
            {"english": "bias", "chinese": "偏见", "learned_date": None, "review_count": 0},
            {"english": "capable", "chinese": "有能力的", "learned_date": None, "review_count": 0},
            {"english": "capacity", "chinese": "能力、容量", "learned_date": None, "review_count": 0},
            {"english": "category", "chinese": "类别", "learned_date": None, "review_count": 0},
            {"english": "challenge", "chinese": "挑战", "learned_date": None, "review_count": 0},
            {"english": "circumstance", "chinese": "情况、环境", "learned_date": None, "review_count": 0},
            {"english": "cite", "chinese": "引用", "learned_date": None, "review_count": 0},
            {"english": "colleague", "chinese": "同事", "learned_date": None, "review_count": 0},
            {"english": "commit", "chinese": "承诺、犯罪", "learned_date": None, "review_count": 0},
            {"english": "communicate", "chinese": "交流", "learned_date": None, "review_count": 0},
            {"english": "community", "chinese": "社区", "learned_date": None, "review_count": 0},
        ]
    
    def save_words(self, words):
        """保存词汇库"""
        with open(self.words_file, 'w', encoding='utf-8') as f:
            json.dump(words, f, ensure_ascii=False, indent=2)
    
    def load_wrong_words(self):
        """加载错题本"""
        if self.wrong_file.exists():
            with open(self.wrong_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        return []
    
    def save_wrong_words(self):
        """保存错题本"""
        with open(self.wrong_file, 'w', encoding='utf-8') as f:
            json.dump(self.wrong_words, f, ensure_ascii=False, indent=2)
    
    def load_learning_log(self):
        """加载学习记录"""
        if self.log_file.exists():
            with open(self.log_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        return {"total_tests": 0, "correct_count": 0, "last_study_date": None}
    
    def save_learning_log(self):
        """保存学习记录"""
        with open(self.log_file, 'w', encoding='utf-8') as f:
            json.dump(self.learning_log, f, ensure_ascii=False, indent=2)
    
    def get_review_words(self):
        """获取需要复习的单词(艾宾浩斯曲线)"""
        today = datetime.now().date()
        review_intervals = [1, 2, 4, 7, 15]  # 艾宾浩斯复习间隔
        
        review_words = []
        for word in self.words:
            if word['learned_date']:
                learned_date = datetime.strptime(word['learned_date'], '%Y-%m-%d').date()
                days_passed = (today - learned_date).days
                
                # 检查是否在复习周期内
                if days_passed in review_intervals:
                    review_words.append(word)
        
        return review_words
    
    def start_quiz(self, mode='random', count=10):
        """开始测试"""
        print("\n" + "="*60)
        print("词汇测试模式")
        print("="*60 + "\n")
        
        # 选择测试词汇
        if mode == 'review':
            test_words = self.get_review_words()
            if not test_words:
                print("[提示] 今天没有需要复习的单词!")
                return
            print(f"今天需要复习 {len(test_words)} 个单词")
        elif mode == 'wrong':
            if not self.wrong_words:
                print("错题本是空的! 你真棒!")
                return
            test_words = self.wrong_words[:count]
            print(f"从错题本中选择 {len(test_words)} 个单词")
        else:
            test_words = random.sample(self.words, min(count, len(self.words)))
            print(f"随机抽取 {len(test_words)} 个单词")
        
        # 开始测试
        correct = 0
        wrong_in_this_test = []
        
        for i, word in enumerate(test_words, 1):
            print(f"\n【{i}/{len(test_words)}】")
            print(f"英文: {word['english']}")
            
            user_answer = input("中文释义: ").strip()
            
            if user_answer.lower() == word['chinese'].lower() or user_answer in word['chinese']:
                print("[v] 正确!")
                correct += 1
                
                # 更新学习日期
                if not word['learned_date']:
                    word['learned_date'] = datetime.now().strftime('%Y-%m-%d')
                word['review_count'] += 1
            else:
                print(f"[x] 错误! 正确答案: {word['chinese']}")
                wrong_in_this_test.append(word)
                
                # 添加到错题本
                if word not in self.wrong_words:
                    self.wrong_words.append(word)
        
        # 统计结果
        accuracy = (correct / len(test_words)) * 100
        print("\n" + "="*60)
        print("测试结果")
        print("="*60)
        print(f"总题数: {len(test_words)}")
        print(f"正确: {correct}")
        print(f"错误: {len(test_words) - correct}")
        print(f"正确率: {accuracy:.1f}%")
        
        # 更新学习记录
        self.learning_log['total_tests'] += len(test_words)
        self.learning_log['correct_count'] += correct
        self.learning_log['last_study_date'] = datetime.now().strftime('%Y-%m-%d')
        
        # 保存数据
        self.save_words(self.words)
        self.save_wrong_words()
        self.save_learning_log()
        
        print("\n数据已保存!")
    
    def search_word(self, keyword):
        """搜索单词"""
        print(f"\n搜索: '{keyword}'")
        print("="*60 + "\n")
        
        results = []
        for word in self.words:
            if keyword.lower() in word['english'].lower() or keyword in word['chinese']:
                results.append(word)
        
        if results:
            for word in results:
                print(f"  {word['english']} - {word['chinese']}")
                if word['learned_date']:
                    print(f"   学习日期: {word['learned_date']}, 复习次数: {word['review_count']}")
                print()
        else:
            print("[未找到] 未找到匹配的单词")
    
    def show_statistics(self):
        """显示学习统计"""
        print("\n" + "="*60)
        print("学习统计")
        print("="*60 + "\n")
        
        # 已学单词数
        learned_count = sum(1 for w in self.words if w['learned_date'])
        
        # 总正确率
        if self.learning_log['total_tests'] > 0:
            total_accuracy = (self.learning_log['correct_count'] / self.learning_log['total_tests']) * 100
        else:
            total_accuracy = 0
        
        # 需要复习的单词
        review_count = len(self.get_review_words())
        
        # 错题本单词数
        wrong_count = len(self.wrong_words)
        
        print(f"词汇库总数: {len(self.words)}")
        print(f"已学单词: {learned_count}")
        print(f"学习进度: {(learned_count/len(self.words)*100):.1f}%")
        print(f"总测试题数: {self.learning_log['total_tests']}")
        print(f"总正确率: {total_accuracy:.1f}%")
        print(f"今日需复习: {review_count} 个")
        print(f"错题本: {wrong_count} 个")
        
        if self.learning_log['last_study_date']:
            last_date = datetime.strptime(self.learning_log['last_study_date'], '%Y-%m-%d').date()
            days_since = (datetime.now().date() - last_date).days
            print(f"最后学习: {self.learning_log['last_study_date']} ({days_since}天前)")
        else:
            print(f"最后学习: 尚未开始")
        
        print()
    
    def print_menu(self):
        """打印菜单"""
        menu = """
┌─────────────────────────────────────┐
│      词汇学习工具菜单               │
├─────────────────────────────────────┤
│  1. 随机测试（10个）                │
│  2. 艾宾浩斯复习提醒                │
│  3. 错题本练习                      │
│  4. 搜索单词                        │
│  5. 学习统计                        │
│  6. 添加单词                        │
│  0. 返回主菜单                      │
└─────────────────────────────────────┘
"""
        print(menu)
    
    def add_word(self):
        """添加单词"""
        print("\n" + "="*60)
        print("添加新单词")
        print("="*60 + "\n")
        
        english = input("英文单词: ").strip()
        if not english:
            print("[错误] 英文单词不能为空")
            return
        
        # 检查是否已存在
        for word in self.words:
            if word['english'].lower() == english.lower():
                print(f"[警告] 单词 '{english}' 已存在: {word['chinese']}")
                return
        
        chinese = input("中文释义: ").strip()
        if not chinese:
            print("[错误] 中文释义不能为空")
            return
        
        new_word = {
            "english": english,
            "chinese": chinese,
            "learned_date": None,
            "review_count": 0
        }
        
        self.words.append(new_word)
        self.save_words(self.words)
        
        print(f"\n[v] 成功添加: {english} - {chinese}")
    
    def run(self):
        """运行交互式界面"""
        while True:
            self.print_menu()
            
            choice = input("请选择功能: ").strip()
            
            if choice == '1':
                self.start_quiz(mode='random', count=10)
            elif choice == '2':
                self.start_quiz(mode='review')
            elif choice == '3':
                self.start_quiz(mode='wrong', count=10)
            elif choice == '4':
                keyword = input("\n请输入搜索关键词: ").strip()
                if keyword:
                    self.search_word(keyword)
            elif choice == '5':
                self.show_statistics()
            elif choice == '6':
                self.add_word()
            elif choice == '0':
                print("\n感谢使用词汇学习工具!")
                break
            else:
                print("\n[警告] 无效选择,请重试")
            
            if choice != '0':
                input("\n按 Enter 键继续...")


def main():
    """主函数"""
    print("""
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║          考研英语词汇学习工具                             ║
║          Vocabulary Learning Tool                         ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
""")
    
    tool = VocabularyTool()
    tool.run()


if __name__ == "__main__":
    main()
