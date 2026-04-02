import javax.swing.*;
import java.awt.*;
import java.awt.event.KeyAdapter;
import java.awt.event.KeyEvent;
import java.io.IOException;
import java.io.OutputStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.util.Objects;

public class client{
    public static String Host;
    public client() {
    }
    public void creatui(){
        //创建一个通讯实体
        clientfun sd=new clientfun();
        //创建连接页面
        JFrame client =creat();
        JTextField host=new JTextField(10);
        JButton hostButton = new JButton("建立通信");
        JPanel panel1=new JPanel();
        JLabel label1=new JLabel("地址");
        panel1.add(label1);
        panel1.add(host);
        panel1.add(hostButton);
        client.add(panel1, BorderLayout.CENTER);
        client.setVisible(true);
        //创建聊天页面
        JFrame clientmain=creat();
        JTextArea show1=new JTextArea(24,64);
        JScrollPane pane1=new JScrollPane(show1);
        show1.setEditable(false);
        JTextField sendMessageField1 = new JTextField(20);
        JPanel panel2=new JPanel();
        JLabel label2=new JLabel("聊天消息：");
        JButton sendButton2 = new JButton("发送");
        JButton endCommunicationButton1 = new JButton("结束通信");
        panel2.add(label2);
        panel2.add(sendMessageField1);
        panel2.add(sendButton2);
        panel2.add(endCommunicationButton1);
        clientmain.add(pane1,BorderLayout.PAGE_START);
        clientmain.add(panel2,BorderLayout.PAGE_END);
        //接收第一次返回消息
        String  re;
        boolean flag=false;
        //监听输入的地址并执行连接
        hostButton.addActionListener(e -> {
            Timer timer = new Timer(3000, e1 -> JOptionPane.showMessageDialog(null, "连接失败：超时", "错误", JOptionPane.ERROR_MESSAGE));
            timer.setRepeats(false);
            timer.start();

            try {
                Host=host.getText();
                sd.send("连接成功", Host);
                timer.stop();
            } catch (IOException ex) {
                throw new RuntimeException(ex);
            }
            String re1 = sd.receive();
            if(!re1.isEmpty()){
                clientmain.setVisible(true);
                client.setVisible(false);
            }
        });
        re = sd.receive();
        if(re.isEmpty()){
            //判断回信为空
            flag=true;
        }
        if(flag){
            //执行两次循环，第一次再去连接服务器，让服务器回显消息提示连接成功并启动接收消息的线程，第二次启动发送的线程
            for(int i=0;i<2;i++){
                try {
                    if(clientfun.i==1){
                        sd.send("连接成功",Host);
                        ReceiveMessageThread receiveMessageThread = new ReceiveMessageThread(sd, show1);
                        receiveMessageThread.start();
                    }
                    else{
                        sendMessageField1.addKeyListener(new EnterKeyListener(sd, Host, sendMessageField1));
                        sendButton2.addActionListener(e -> {
                            SendMessageThread sendMessageThread = new SendMessageThread(sd, Host, sendMessageField1);
                            sendMessageThread.start();
                        });
                        //清空监听次数
                        if(clientfun.i!=2){
                            sendButton2.removeActionListener(sendButton2.getActionListeners()[0]);
                        }
                    }
                    endCommunicationButton1.addActionListener(e -> {
                        end endconnect=new end(clientmain,client);
                        endconnect.run();
                    });
                } catch (IOException ex) {
                    throw new RuntimeException(ex);
                }
            }
        }
        //阻塞主线程结束
        while(true){}
    }
    public JFrame creat(){
        JFrame frame = new JFrame("通讯");
        frame.setSize(800, 600);
        ImageIcon coin=new ImageIcon("微信图片_20240419142656.png");
        frame.setIconImage(coin.getImage());
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        return frame;
    }
}
class clientfun {
    //创建一个信号量检查是否是第一次做连接
    public static int i=0;
    public void send(String messagesend,String host) throws IOException {
        //开启一个连接选定服务地址的5000端口
        Socket sc = new Socket(host, 500);
        //打开输出流
        OutputStream os = sc.getOutputStream();
        //将发送的消息以utf-8编码的字节发送
        os.write(messagesend.getBytes(StandardCharsets.UTF_8));
        //关闭输出流
        os.close();
        //关闭连接
        sc.close();
        i++;
    }
    public String receive() {
        try {
            // 创建一个DatagramSocket对象，绑定到指定的端口
            DatagramSocket socket = new DatagramSocket(1000);
            // 创建一个DatagramPacket对象，用于接收数据包
            byte[] buffer = new byte[1024];
            DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
            // 使用DatagramSocket对象的receive()方法接收数据包
            socket.receive(packet);
            // 从DatagramPacket对象中获取接收到的消息
            String message = new String(packet.getData(), 0, packet.getLength(), StandardCharsets.UTF_8);
            // 关闭DatagramSocket对象
            socket.close();
            //回传收到的消息
            return message;
        } catch (IOException e) {
            e.printStackTrace();
            return "连接异常，请结束重新连接！";
        }
    }
}
//创建一个专门接收消息的线程，防止主线程阻塞
class ReceiveMessageThread extends Thread {
    private final clientfun sd;
    private final JTextArea show1;

    public ReceiveMessageThread(clientfun sd, JTextArea show1) {
        this.sd = sd;
        this.show1 = show1;
    }

    @Override
    public void run() {
        while (true) {
            String message = sd.receive();
            if(Objects.equals(message, "")){
                show1.append(message);
            }else{
                show1.append(message+"\n");
            }
        }
    }
}
//创建一个发送消息的线程，将它绑定到发送按钮上点击一次就实现一次，其余时间阻塞
class SendMessageThread extends Thread {
    private final clientfun sd;
    private final String host;
    private final JTextField sendMessageField1;
    public SendMessageThread(clientfun sd, String host, JTextField sendMessageField1) {
        this.sd = sd;
        this.host = host;
        this.sendMessageField1 = sendMessageField1;
    }
    @Override
    public void run() {
        try {
            String mess = sendMessageField1.getText();
            sendMessageField1.setText("");
            if (mess.equals("连接异常，请结束重新连接！")) {
                JOptionPane.showMessageDialog(null, "连接异常，请结束重新连接！", "错误", JOptionPane.ERROR_MESSAGE);
                System.exit(0);
            }
            sd.send(mess, host);
        } catch (IOException ex) {
            throw new RuntimeException(ex);
        }
    }
}
class end implements Runnable{
    private final JFrame clientmain;
    private final JFrame client;
    public end(JFrame clientmain, JFrame client) {
        this.clientmain=clientmain;
        this.client=client;
    }
    @Override
    public void run() {
        clientmain.setVisible(false);
        client.setVisible(true);
    }
}
class EnterKeyListener extends KeyAdapter {
    private final clientfun sd;
    private final String host;
    private final JTextField sendMessageField1;

    public EnterKeyListener(clientfun sd, String host, JTextField sendMessageField1) {
        this.sd = sd;
        this.host = host;
        this.sendMessageField1 = sendMessageField1;
    }

    @Override
    public void keyPressed(KeyEvent e) {
        if (e.getKeyCode() == KeyEvent.VK_ENTER) {
            try {
                String mess = sendMessageField1.getText();
                sendMessageField1.setText("");
                if (mess.equals("连接异常，请结束重新连接！")) {
                    JOptionPane.showMessageDialog(null, "连接异常，请结束重新连接！", "错误", JOptionPane.ERROR_MESSAGE);
                    System.exit(0);
                }
                sd.send(mess, host);
            } catch (IOException ex) {
                throw new RuntimeException(ex);
            }
        }
    }
}

