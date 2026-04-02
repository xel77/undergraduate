import os
import json
import uuid
import hashlib
from datetime import datetime, timedelta
import tkinter as tk
from tkinter import messagebox, simpledialog

CONFIG_FILE = "license_info.json"

# 获取机器唯一标识
def get_machine_id():
    return str(uuid.getnode())

# 使用 MD5 生成注册码
def generate_registration_code(machine_id):
    return hashlib.md5(machine_id.encode()).hexdigest()

# 加载本地许可信息（首次生成）
def load_license():
    if os.path.exists(CONFIG_FILE):
        with open(CONFIG_FILE, "r") as f:
            return json.load(f)
    else:
        return {
            "start_date": datetime.now().strftime("%Y-%m-%d"),
            "usage_count": 0,
            "registered": False
        }

# 保存许可信息
def save_license(data):
    with open(CONFIG_FILE, "w") as f:
        json.dump(data, f)

# 主应用类
class LicenseApp:
    def __init__(self, root):
        self.root = root
        self.root.title("软件版权保护实验")
        self.license = load_license()

        self.machine_id = get_machine_id()
        self.expected_code = generate_registration_code(self.machine_id)

        self.check_license()

        self.build_gui()

    def check_license(self):
        # 检查试用期
        start_date = datetime.strptime(self.license["start_date"], "%Y-%m-%d")
        if not self.license["registered"] and datetime.now() > start_date + timedelta(days=7):
            messagebox.showerror("试用过期", "试用已超过 7 天，请注册后继续使用。")
            exit()

        # 检查试用次数
        if not self.license["registered"] and self.license["usage_count"] >= 10:
            messagebox.showerror("试用次数已满", "已达到最大试用次数（10次），请注册后继续使用。")
            exit()

        # 增加使用次数
        self.license["usage_count"] += 1
        save_license(self.license)

    def build_gui(self):
        tk.Label(self.root, text="欢迎使用本软件").pack(pady=10)

        self.status_label = tk.Label(self.root, text="")
        self.status_label.pack()

        self.basic_button = tk.Button(self.root, text="基础功能", command=self.use_basic)
        self.basic_button.pack(pady=5)

        self.advanced_button = tk.Button(self.root, text="高级功能", command=self.use_advanced)
        self.advanced_button.pack(pady=5)

        self.register_button = tk.Button(self.root, text="注册软件", command=self.register)
        self.register_button.pack(pady=10)

        self.update_ui()

    def update_ui(self):
        if self.license["registered"]:
            self.status_label.config(text="✅ 已注册用户，功能全部开放")
            self.advanced_button.config(state="normal")
        else:
            remaining_days = 7 - (datetime.now() - datetime.strptime(self.license["start_date"], "%Y-%m-%d")).days
            remaining_uses = 10 - self.license["usage_count"]
            self.status_label.config(
                text=f"🔒 未注册用户 | 剩余天数：{remaining_days} | 剩余次数：{remaining_uses}")
            self.advanced_button.config(state="disabled")

    def use_basic(self):
        messagebox.showinfo("基础功能", "这是一个可用的基础功能模块。")

    def use_advanced(self):
        messagebox.showinfo("高级功能", "这是一个已解锁的高级功能模块。")

    def register(self):
        code = simpledialog.askstring("注册", f"请输入注册码（设备码：{self.machine_id}）")
        if code == self.expected_code:
            self.license["registered"] = True
            save_license(self.license)
            messagebox.showinfo("注册成功", "恭喜，注册成功！功能已全部解锁。")
            self.update_ui()
        else:
            messagebox.showerror("注册失败", "注册码错误，请重新输入。")

if __name__ == "__main__":
    root = tk.Tk()
    app = LicenseApp(root)
    root.mainloop()
