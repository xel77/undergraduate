import os
import pandas as pd
from collections import defaultdict

# 定义泛化函数
def generalize_age(age, level):
    """
    年龄的泛化函数。

    Parameters:
        age (int or str): 具体年龄或已泛化的年龄范围。
        level (int): 泛化层级。

    Returns:
        str/int: 泛化后的年龄。
    """
    if level == 0:
        return age
    elif level == 1:
        if isinstance(age, int):
            return f"{(age // 10) * 10}-{(age // 10) * 10 + 9}"  # 按10岁区间泛化
        else:
            return age  # 已经泛化，不再变化
    elif level == 2:
        if isinstance(age, int):
            return f"{(age // 5) * 5}-{(age // 5) * 5 + 4}"    # 按5岁区间泛化
        elif isinstance(age, str) and age != '*':
            return '*'  # 进一步泛化为 '*'
        else:
            return '*'  # 已经完全泛化
    else:
        return '*'  # 默认完全泛化

def generalize_education(education, level):
    """
    教育程度的泛化函数。

    Parameters:
        education (str): 具体教育程度。
        level (int): 泛化层级。

    Returns:
        str: 泛化后的教育程度。
    """
    if level == 0:
        return education
    elif level == 1:
        high_school = [
            'Preschool', '1st-4th', '5th-6th', '7th-8th',
            '9th', '10th', '11th', '12th', 'HS-grad'
        ]
        college = ['Some-college', 'Assoc-acdm', 'Assoc-voc']
        advanced = ['Bachelors', 'Masters', 'Doctorate', 'Prof-school']
        if education in high_school:
            return 'High School'
        elif education in college:
            return 'College'
        elif education in advanced:
            return 'Advanced'
        else:
            return 'Other'
    else:
        return '*'  # 完全泛化

def generalize_marital_status(status, level):
    """
    婚姻状况的泛化函数。

    Parameters:
        status (str): 具体婚姻状况。
        level (int): 泛化层级。

    Returns:
        str: 泛化后的婚姻状况。
    """
    if level == 0:
        return status
    elif level == 1:
        married = [
            'Married-civ-spouse', 'Married-spouse-absent', 'Married-AF-spouse'
        ]
        if status in married:
            return 'Married'
        else:
            return status
    else:
        return '*'  # 完全泛化

def generalize_race(race, level):
    """
    种族的泛化函数。

    Parameters:
        race (str): 具体种族。
        level (int): 泛化层级。

    Returns:
        str: 泛化后的种族。
    """
    if level == 0:
        return race
    elif level == 1:
        majority = ['White', 'Black']
        minority = ['Asian-Pac-Islander', 'Amer-Indian-Eskimo', 'Other']
        if race in majority:
            return 'Majority'
        elif race in minority:
            return 'Minority'
        else:
            return 'Other'
    else:
        return '*'  # 完全泛化

# k-匿名性检测函数
def is_k_anonymous(df, quasi_identifiers, k):
    """
    检查数据集是否满足k-匿名性。

    Parameters:
        df (pd.DataFrame): 数据集。
        quasi_identifiers (list): 准标识符属性列表。
        k (int): k值。

    Returns:
        bool: 是否满足k-匿名性。
    """
    groups = df.groupby(quasi_identifiers).size()
    return all(groups >= k)

# DataFly算法实现，带抑制策略
def datafly_anonymization_with_suppression(df, k, quasi_identifiers, generalization_levels, max_generalization):
    """
    实现DataFly算法进行k-匿名化，增加抑制策略。

    Parameters:
        df (pd.DataFrame): 原始数据集。
        k (int): k值。
        quasi_identifiers (list): 准标识符属性列表。
        generalization_levels (dict): 当前泛化层级字典。
        max_generalization (dict): 最大泛化层级字典。

    Returns:
        pd.DataFrame: 匿名化后的数据集。
        dict: 最终的泛化层级。
    """
    df_anonymized = df.copy()
    iteration = 0

    while not is_k_anonymous(df_anonymized, quasi_identifiers, k):
        iteration += 1
        print(f"Iteration {iteration}: Data not yet k-anonymous.")

        # 统计每个准标识符的取值个数
        distinct_counts = {}
        for qi in quasi_identifiers:
            distinct_counts[qi] = df_anonymized[qi].nunique()

        # 选择取值最多的属性进行泛化
        # 只考虑尚未达到最大泛化层次的属性
        available_attrs = [
            qi for qi in quasi_identifiers
            if generalization_levels[qi] < max_generalization[qi]
        ]
        if available_attrs:
            # 在可泛化的属性中选择取值最多的属性
            available_counts = {qi: distinct_counts[qi] for qi in available_attrs}
            attr_to_generalize = max(available_counts, key=available_counts.get)

            # 增加该属性的泛化层次
            generalization_levels[attr_to_generalize] += 1
            current_level = generalization_levels[attr_to_generalize]
            print(f"泛化属性: {attr_to_generalize}, 当前层级: {current_level}")

            # 应用泛化
            if attr_to_generalize == 'age':
                df_anonymized.loc[:, 'age'] = df_anonymized['age'].apply(
                    lambda x: generalize_age(x, current_level)
                )
            elif attr_to_generalize == 'education':
                df_anonymized.loc[:, 'education'] = df_anonymized['education'].apply(
                    lambda x: generalize_education(x, current_level)
                )
            elif attr_to_generalize == 'marital-status':
                df_anonymized.loc[:, 'marital-status'] = df_anonymized['marital-status'].apply(
                    lambda x: generalize_marital_status(x, current_level)
                )
            elif attr_to_generalize == 'race':
                df_anonymized.loc[:, 'race'] = df_anonymized['race'].apply(
                    lambda x: generalize_race(x, current_level)
                )
            else:
                print(f"未知的属性: {attr_to_generalize}")
        else:
            # 所有属性已达到最大泛化层次，进行抑制
            print("所有属性已达到最大泛化层次，开始抑制不满足k-匿名性的记录。")
            # 找到不满足k-匿名性的等价类
            groups = df_anonymized.groupby(quasi_identifiers).size().reset_index(name='count')
            violating_groups = groups[groups['count'] < k]
            if violating_groups.empty:
                break  # 已满足k-匿名性
            # 删除这些不满足k-匿名性的记录
            df_anonymized = df_anonymized.merge(
                violating_groups[quasi_identifiers],
                on=quasi_identifiers,
                how='left',
                indicator=True
            )
            df_anonymized = df_anonymized[df_anonymized['_merge'] == 'left_only']
            df_anonymized.drop(columns=['_merge'], inplace=True)
            print(f"抑制了 {len(violating_groups)} 个等价类中的记录。")

    print("完成匿名化处理。")
    return df_anonymized, generalization_levels

# 计算评价指标
def calculate_precision_distortion(generalization_levels, max_generalization, df_anonymized, quasi_identifiers):
    """
    计算精确率和失真率。

    Parameters:
        generalization_levels (dict): 泛化层级字典。
        max_generalization (dict): 最大泛化层级字典。
        df_anonymized (pd.DataFrame): 匿名化后的数据集。
        quasi_identifiers (list): 准标识符属性列表。

    Returns:
        float: 精确率。
        float: 失真率。
    """
    total_records = len(df_anonymized)
    num_qi = len(quasi_identifiers)

    if total_records == 0:
        return 0.0, 0.0  # 避免除以零

    # 计算每个准标识符的泛化程度
    generalization_degrees = {}
    for qi in quasi_identifiers:
        generalization_degrees[qi] = generalization_levels[qi] / max_generalization[qi]

    # 计算精确率
    # precision=1-（所有记录的各个准标识符的泛化程度之和）/（总记录数*准标识符个数）
    sum_generalization = sum(generalization_degrees.values()) * total_records
    precision = 1 - (sum_generalization) / (total_records * num_qi)

    # 计算失真率
    # distortion=各个准标识符的泛化程度之和/准标识符个数
    distortion = sum(generalization_degrees.values()) / num_qi

    return precision, distortion

# 数据加载函数
def load_adult_data(data_dir='data'):
    """
    加载Adult数据集的训练和测试数据，并进行初步解析和预处理。

    参数:
        data_dir (str): 存放adult.data和adult.test文件的目录路径。默认为当前目录下的'data'文件夹。

    返回:
        pd.DataFrame: 合并并预处理后的Adult数据集。
    """
    # 定义列名
    column_names = [
        'age', 'workclass', 'fnlwgt', 'education', 'education-num',
        'marital-status', 'occupation', 'relationship', 'race', 'sex',
        'capital-gain', 'capital-loss', 'hours-per-week', 'native-country', 'income'
    ]

    # 构建文件路径
    train_data_path = os.path.join(data_dir, 'adult.data')
    test_data_path = os.path.join(data_dir, 'adult.test')

    # 检查文件是否存在
    if not os.path.exists(train_data_path):
        raise FileNotFoundError(f"训练数据文件未找到: {train_data_path}")
    if not os.path.exists(test_data_path):
        raise FileNotFoundError(f"测试数据文件未找到: {test_data_path}")

    # 加载训练数据
    train_data = pd.read_csv(
        train_data_path,
        header=None,
        names=column_names,
        na_values=' ?',         # 将 ' ?' 视为缺失值
        skipinitialspace=True   # 跳过每行开头的空格
    )

    # 加载测试数据
    test_data = pd.read_csv(
        test_data_path,
        header=None,
        names=column_names,
        na_values=' ?',         # 将 ' ?' 视为缺失值
        skiprows=1,             # 跳过第一行说明
        skipinitialspace=True   # 跳过每行开头的空格
    )

    # 合并训练数据和测试数据
    combined_data = pd.concat([train_data, test_data], ignore_index=True)

    # 处理缺失值
    initial_count = len(combined_data)
    combined_data.dropna(inplace=True)
    final_count = len(combined_data)
    print(f"删除缺失值前的记录数：{initial_count}")
    print(f"删除缺失值后的记录数：{final_count}")

    return combined_data

# 主流程
def main():
    # 定义数据目录
    data_directory = 'data'  # 根据实际情况修改路径

    # 加载数据
    try:
        data = load_adult_data(data_dir=data_directory)
    except FileNotFoundError as e:
        print(e)
        return

    # 查看数据基本信息
    print("\n数据集基本信息：")
    print(data.info())

    # 查看前5条记录
    print("\n数据集前5条记录：")
    print(data.head())

    # 检查是否存在缺失值
    missing_values = data.isnull().sum()
    print("\n每列的缺失值数量：")
    print(missing_values)

    # 选择需要的字段
    selected_columns = ['age', 'education', 'marital-status', 'race', 'income']
    data_selected = data[selected_columns].copy()  # 使用 .copy() 避免 SettingWithCopyWarning

    # 重置索引
    data_selected.reset_index(drop=True, inplace=True)

    print("\n预处理后的数据集前5条记录：")
    print(data_selected.head())

    # 确保 age 列是整数类型，如果有问题的值，则进行处理
    data_selected.loc[:, 'age'] = pd.to_numeric(data_selected['age'], errors='coerce')  # 将无法转换的值变为 NaN
    data_selected.dropna(subset=['age'], inplace=True)  # 删除 age 列中存在问题的记录

    # 将 age 转换为整数类型
    data_selected.loc[:, 'age'] = data_selected['age'].astype(int)

    print("\n处理后的数据集前5条记录（确保age为整数）：")
    print(data_selected.head())

    # 定义k值
    k = 5  # 可以根据需求调整

    # 定义准标识符属性
    quasi_identifiers = ['age', 'education', 'marital-status', 'race']

    # 定义泛化层次字典，记录每个属性当前的泛化层级
    generalization_levels = {qi: 0 for qi in quasi_identifiers}

    # 定义每个属性的最大泛化层级
    max_generalization_levels = {
        'age': 3,              # 0: exact, 1: 10s, 2: 5s, 3: '*'
        'education': 2,        # 0: exact, 1: grouped, 2: '*'
        'marital-status': 2,   # 0: exact, 1: grouped, 2: '*'
        'race': 2              # 0: exact, 1: grouped, 2: '*'
    }

    # 应用DataFly算法（带抑制策略）
    anonymized_data, final_generalization_levels = datafly_anonymization_with_suppression(
        data_selected,
        k,
        quasi_identifiers,
        generalization_levels,
        max_generalization_levels
    )

    print("\n匿名化后的数据示例：")
    print(anonymized_data.head())

    print("\n最终的泛化层次：")
    print(final_generalization_levels)

    # 计算评价指标
    precision, distortion = calculate_precision_distortion(
        final_generalization_levels,
        max_generalization_levels,
        anonymized_data,
        quasi_identifiers
    )

    print(f"\n精确率（Precision）：{precision:.4f}")
    print(f"失真率（Distortion）：{distortion:.4f}")

    # 保存匿名化后的数据（可选）
    output_path = os.path.join(data_directory, 'anonymized_adult_data.csv')
    anonymized_data.to_csv(output_path, index=False)
    print(f"\n匿名化后的数据已保存到 {output_path}")

if __name__ == "__main__":
    main()
