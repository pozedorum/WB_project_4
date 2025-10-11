#!/usr/bin/env python3
# create_test_file.py

import random

def create_test_file():
    filename = "test_large_file.txt"
    total_lines = 300000
    test_word = "test"
    test_count = 20
    
    print(f"Creating test file with {total_lines} lines and {test_count} occurrences of '{test_word}'...")
    
    # Создаем список всех строк
    lines = []
    test_index = 0
    every = total_lines/test_count
    # Добавляем строки
    for i in range(total_lines):
        lines.append(f"This is a normal line without keywords. Line number: {i+1}\n")
        if i == test_index * every:
            lines.append(f"This is a line with the {test_word} word. Line number: {i+1}\n")
            test_index+=1
    
    # Перемешиваем строки (опционально)
    # random.shuffle(lines)
    
    # Записываем в файл
    with open(filename, 'w') as f:
        f.writelines(lines)
    
    print(f"File '{filename}' created successfully!")
    print(f"Lines with '{test_word}': {test_count}")
    print(f"Total lines: {total_lines}")

if __name__ == "__main__":
    create_test_file()