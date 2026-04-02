// datafly.ts

import * as fs from 'fs';
import * as csv from 'csv-parser';

interface RecordData {
    [key: string]: string;
}

class Datafly {
    data: RecordData[];
    quasiIdentifiers: string[];
    k: number;
    generalizationLevels: { [key: string]: number };
    maxGeneralizationLevels: { [key: string]: number };

    constructor(data: RecordData[], quasiIdentifiers: string[], k: number) {
        this.data = data;
        this.quasiIdentifiers = quasiIdentifiers;
        this.k = k;
        // 初始化每个准标识符的一般化级别
        this.generalizationLevels = {};
        this.maxGeneralizationLevels = {};
        this.quasiIdentifiers.forEach(col => {
            this.generalizationLevels[col] = 0;
            this.maxGeneralizationLevels[col] = this.getMaxGeneralizationLevel(col);
        });
    }

    // 获取每个准标识符的最大一般化级别
    getMaxGeneralizationLevel(col: string): number {
        const maxLevels: { [key: string]: number } = {
            '年龄': 2,
            '区号': 3,
            '健康状况': 2,
            '性别': 1,
            // 添加更多列的最大一般化级别
        };
        return maxLevels[col] || 1;
    }

    // 一般化策略
    generalize(record: RecordData): RecordData {
        const generalizedRecord: RecordData = { ...record };

        this.quasiIdentifiers.forEach((col) => {
            if (col === '年龄') {
                generalizedRecord[col] = this.generalizeAge(record[col], this.generalizationLevels[col]);
            } else if (col === '区号') {
                generalizedRecord[col] = this.generalizeAreaCode(record[col], this.generalizationLevels[col]);
            } else if (col === '健康状况') {
                generalizedRecord[col] = this.generalizeHealthStatus(record[col], this.generalizationLevels[col]);
            } else if (col === '性别') {
                generalizedRecord[col] = this.generalizeGender(record[col], this.generalizationLevels[col]);
            }
            // 可以根据需要添加更多列的一般化策略
        });

        return generalizedRecord;
    }

    // 年龄的一般化策略
    generalizeAge(age: string, level: number): string {
        if (level === 0) {
            const ageNum = parseInt(age, 10);
            if (isNaN(ageNum)) {
                // 处理描述性年龄
                switch (age) {
                    case '少':
                        return '<20';
                    case '中':
                        return '20-59';
                    case '老':
                        return '60+';
                    default:
                        return '未知';
                }
            } else {
                if (ageNum < 20) {
                    return '<20';
                } else if (ageNum < 40) {
                    return '20-39';
                } else if (ageNum < 60) {
                    return '40-59';
                } else {
                    return '60+';
                }
            }
        } else if (level === 1) {
            // 更广泛的年龄段
            const ageNum = parseInt(age, 10);
            if (isNaN(ageNum)) {
                // 基于描述性年龄
                switch (age) {
                    case '少':
                    case '中':
                        return '<40';
                    case '老':
                        return '40+';
                    default:
                        return '未知';
                }
            } else {
                if (ageNum < 40) {
                    return '<40';
                } else {
                    return '40+';
                }
            }
        } else {
            // 最广泛的类别
            return '*';
        }
    }

    // 区号的一般化策略
    generalizeAreaCode(areaCode: string, level: number): string {
        if (level === 0) {
            // Level 0: 0595 -> 059XX
            return areaCode.substring(0, 3) + 'XX';
        } else if (level === 1) {
            // Level 1: 0595 -> 05XX
            return areaCode.substring(0, 2) + 'XX';
        } else if (level === 2) {
            // Level 2: 0595 -> 0XXX
            return areaCode.substring(0, 1) + 'XXX';
        } else {
            // 最广泛的类别
            return '*';
        }
    }

    // 健康状况的一般化策略
    generalizeHealthStatus(status: string, level: number): string {
        if (level === 0) {
            const generalCategories: { [key: string]: string } = {
                '发烧': '呼吸系统疾病',
                '感冒': '呼吸系统疾病',
                '肺炎': '呼吸系统疾病',
                '高血压': '心血管疾病',
                '糖尿病': '代谢疾病',
                '脑溢血': '神经系统疾病',
                '心脏病': '心血管疾病',
                '癌症': '肿瘤',
                // 添加更多具体疾病的归类
            };
            return generalCategories[status] || '其他疾病';
        } else if (level === 1) {
            // 更广泛的类别
            const broaderCategories: { [key: string]: string } = {
                '呼吸系统疾病': '疾病',
                '心血管疾病': '疾病',
                '代谢疾病': '疾病',
                '神经系统疾病': '疾病',
                '肿瘤': '疾病',
                '其他疾病': '疾病',
                // 所有类别归为 '疾病'
            };
            const categoryLevel1 = this.generalizeHealthStatus(status, 0);
            return broaderCategories[categoryLevel1] || '疾病';
        } else {
            // 最广泛的类别
            return '*';
        }
    }

    // 性别的一般化策略
    generalizeGender(gender: string, level: number): string {
        if (level === 0) {
            return gender;
        } else {
            return '不明确';
        }
    }

    // 计算频率
    computeFrequency(data: RecordData[]): Map<string, number> {
        const freqMap = new Map<string, number>();
        data.forEach((record) => {
            const key = this.quasiIdentifiers.map((col) => record[col]).join('|');
            freqMap.set(key, (freqMap.get(key) || 0) + 1);
        });
        return freqMap;
    }

    // 进行一般化
    applyGeneralization(): void {
        this.data = this.data.map((record) => this.generalize(record));

        // 增加每个准标识符的一般化级别
        this.quasiIdentifiers.forEach(col => {
            if (this.generalizationLevels[col] < this.maxGeneralizationLevels[col]) {
                this.generalizationLevels[col] += 1;
            }
        });
    }

    // 进行抑制（删除低频记录）
    suppress(lowFreqKeys: Set<string>): void {
        this.data = this.data.filter((record) => {
            const key = this.quasiIdentifiers.map((col) => record[col]).join('|');
            return !lowFreqKeys.has(key);
        });
    }

    // 运行 Datafly 算法
    run(): RecordData[] {
        let iteration = 0;
        while (true) {
            iteration += 1;
            console.log(`\n--- Iteration ${iteration} ---`);
            const freqMap = this.computeFrequency(this.data);
            const lowFreqKeys = new Set<string>();

            freqMap.forEach((count, key) => {
                if (count < this.k) {
                    lowFreqKeys.add(key);
                }
            });

            console.log(`Low frequency combinations: ${lowFreqKeys.size}`);

            if (lowFreqKeys.size === 0) {
                console.log('All combinations meet k-anonymity.');
                break; // 满足 k-匿名性
            }

            // 检查是否可以进一步一般化
            const canGeneralize = this.quasiIdentifiers.some(col => this.generalizationLevels[col] < this.maxGeneralizationLevels[col]);

            if (canGeneralize) {
                // 一般化
                console.log('Applying generalization...');
                this.applyGeneralization();
            } else {
                // 无法进一步一般化，进行抑制
                console.log('Cannot generalize further. Suppressing low frequency records...');
                this.suppress(lowFreqKeys);
                break;
            }
        }

        return this.data;
    }
}

function getData1(): RecordData[] {
    const data1: RecordData[] = [
        { ID: '1', 性别: '男', 年龄: '中', 区号: '0595', 健康状况: '发烧' },
        { ID: '2', 性别: '男', 年龄: '老', 区号: '0592', 健康状况: '感冒' },
        { ID: '3', 性别: '女', 年龄: '少', 区号: '0594', 健康状况: '肺炎' },
        { ID: '4', 性别: '女', 年龄: '中', 区号: '0591', 健康状况: '高血压' },
        { ID: '5', 性别: '男', 年龄: '少', 区号: '0596', 健康状况: '糖尿病' },
        { ID: '6', 性别: '男', 年龄: '老', 区号: '0593', 健康状况: '脑溢血' },
        { ID: '7', 性别: '女', 年龄: '少', 区号: '0594', 健康状况: '心脏病' },
        { ID: '8', 性别: '男', 年龄: '中', 区号: '0593', 健康状况: '发烧' },
        { ID: '9', 性别: '女', 年龄: '少', 区号: '0594', 健康状况: '感冒' },
        { ID: '10', 性别: '女', 年龄: '中', 区号: '0592', 健康状况: '肺炎' },
        { ID: '11', 性别: '女', 年龄: '老', 区号: '0594', 健康状况: '癌症' },
        { ID: '12', 性别: '男', 年龄: '老', 区号: '0596', 健康状况: '肺炎' },
        { ID: '13', 性别: '女', 年龄: '中', 区号: '0596', 健康状况: '高血压' },
        { ID: '14', 性别: '女', 年龄: '老', 区号: '0593', 健康状况: '高血压' },
        { ID: '15', 性别: '男', 年龄: '中', 区号: '0595', 健康状况: '糖尿病' },
        { ID: '16', 性别: '男', 年龄: '少', 区号: '0592', 健康状况: '脑溢血' },
        { ID: '17', 性别: '女', 年龄: '少', 区号: '0594', 健康状况: '肺炎' },
        { ID: '18', 性别: '女', 年龄: '中', 区号: '0595', 健康状况: '感冒' },
        { ID: '19', 性别: '男', 年龄: '少', 区号: '0592', 健康状况: '肺炎' },
        { ID: '20', 性别: '女', 年龄: '少', 区号: '0596', 健康状况: '发烧' },
        { ID: '21', 性别: '男', 年龄: '中', 区号: '0594', 健康状况: '感冒' },
        { ID: '22', 性别: '男', 年龄: '老', 区号: '0593', 健康状况: '肺炎' },
        { ID: '23', 性别: '女', 年龄: '中', 区号: '0591', 健康状况: '糖尿病' },
        { ID: '24', 性别: '男', 年龄: '少', 区号: '0595', 健康状况: '脑溢血' },
        { ID: '25', 性别: '男', 年龄: '中', 区号: '0593', 健康状况: '感冒' },
        { ID: '26', 性别: '男', 年龄: '老', 区号: '0595', 健康状况: '肺炎' },
        { ID: '27', 性别: '女', 年龄: '少', 区号: '0592', 健康状况: '发烧' },
        { ID: '28', 性别: '女', 年龄: '老', 区号: '0594', 健康状况: '肺炎' },
        { ID: '29', 性别: '男', 年龄: '中', 区号: '0591', 健康状况: '糖尿病' },
        { ID: '30', 性别: '女', 年龄: '老', 区号: '0593', 健康状况: '发烧' },
    ];
    return data1;
}

// 主函数：应用 Datafly 算法
async function main() {
    // 测试数据一
    const data1 = getData1();
    // 将 '健康状况' 也作为准标识符
    const quasiIdentifiers1 = ['性别', '年龄', '区号', '健康状况'];
    const kValues = [2, 3];

    for (const k of kValues) {
        console.log(`\n=== Datafly Algorithm with k=${k} on Data1 ===`);
        const datafly1 = new Datafly(data1, quasiIdentifiers1, k);
        const anonymizedData1 = datafly1.run();
        console.table(anonymizedData1);
    }
}

main().catch((err) => {
    console.error(err);
});
