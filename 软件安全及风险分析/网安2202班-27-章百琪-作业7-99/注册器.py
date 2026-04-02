import uuid
import hashlib
import tkinter as tk
from tkinter import messagebox
def get_machine_id():
    return str(uuid.getnode())
def generate_registration_code(machine_id):
    return hashlib.md5(machine_id.encode()).hexdigest()
def generate_and_show():
    try:
        machine_id = get_machine_id()
        reg_code = generate_registration_code(machine_id)
        machine_entry.delete(0, tk.END)
        machine_entry.insert(0, machine_id)
        reg_entry.delete(0, tk.END)
        reg_entry.insert(0, reg_code)
    except Exception as e:
        messagebox.showerror("错误", f"生成失败: {str(e)}")
# 创建主窗口
root = tk.Tk()
root.title("机器码生成器")
root.geometry("400x200")
# 机器码部分
tk.Label(root, text="设备码:").pack(pady=(10, 0))
machine_entry = tk.Entry(root, width=40)
machine_entry.pack()
# 注册码部分
tk.Label(root, text="注册码:").pack(pady=(10, 0))
reg_entry = tk.Entry(root, width=40)
reg_entry.pack()
# 生成按钮
tk.Button(root, text="生成注册码", command=generate_and_show).pack(pady=20)
root.mainloop()