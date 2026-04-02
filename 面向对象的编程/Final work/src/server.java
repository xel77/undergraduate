import javax.swing.*;
import java.awt.*;
import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.Objects;

public class server {
    public static int i=0;
    public static  JFrame mainui;
    //构造无参函数
    public server(JFrame mainui) {
        server.mainui =mainui;
    }
    //创建窗体
    public JFrame creat(){
        JFrame frame = new JFrame("服务器");
        frame.setSize(800, 600);
        frame.setLayout(new FlowLayout());
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        return frame;
    }
    public void creatui() throws IOException{
        boolean flag=false;
        InetAddress ip=InetAddress.getLocalHost();
        String ipAddress = "您的ip为:"+ip.getHostAddress();
        //建立连接
        accept sc=new accept();
        //创建主体页面
        JFrame servermain=creat();
        JTextArea show1=new JTextArea(24,64);
        JScrollPane pane1=new JScrollPane(show1);
        JLabel lable=new JLabel(ipAddress);
        show1.setEditable(false);
        JPanel panel2=new JPanel();
        JButton endCommunicationButton1 = new JButton("结束通信");
        panel2.add(lable);
        panel2.add(endCommunicationButton1);
        servermain.add(pane1,BorderLayout.PAGE_START);
        servermain.add(panel2,BorderLayout.PAGE_END);
        servermain.setVisible(true);
        endCommunicationButton1.addActionListener(e -> {
            endserver endconnect=new endserver(servermain,mainui);
            endconnect.run();
        });
        String mess=sc.connect();
        if(Objects.equals(mess, "")){
            flag = true;
            sc.send(mess);
        }
        while(flag){
            String mes=sc.connect();
            if(Objects.equals(mes, "")){
                show1.append(mes);
            }else{
                show1.append(mes+"\n");
            }
            if(i==0){
                sc.send(mes);
            }else{
                show1.append(mes +"\n");
                sc.send(mes);
                i++;

            }
        }
    }

    public static void main(String[] args) throws IOException {
        server server=new server(mainui);
        server.creatui();
    }
}
class accept {
    public static String[] iparry=new String[1024];
    public static int k=-1;
    public static String ip=null;
    public static int p=0;
    public static int t=0;
    //打开端口5000创建tcp连接
    ServerSocket sc= new ServerSocket(500);
    public accept() throws IOException {
    }
    public String connect() throws IOException {
        int flag=1;
        //建立连接
        Socket ce = sc.accept();
        //获得字节流
        InputStream is = ce.getInputStream();
        //转化成字符流
        InputStreamReader isr = new InputStreamReader(is, StandardCharsets.UTF_8);
        //设置成缓冲流
        BufferedReader br = new BufferedReader(isr);
        ip = ce.getInetAddress().getHostAddress();
        for(int j=0;j<k;j++){
            if(ip.equals(iparry[j])){
                flag=0;
                break;
            }
        }
        if(flag==1){
            t=k;
            iparry[++k]=ip;
            p=0;
        }
        //输出消息
        String message = br.readLine();
        //关闭通道
        is.close();
        //关闭连接
        ce.close();
        if(p==1){
            message=ip+"连接成功";
        }
        else if(p==0){
            message="";
        }else{
            message="来自"+ip+"的消息:"+message;
        }
        p++;
        return message;
    }
    public void send(String in) throws IOException{
        byte[] bytes = in.getBytes(StandardCharsets.UTF_8);
        DatagramSocket ds = new DatagramSocket();
        for (int j = 0; j < k + 1; j++) {
            DatagramPacket dp = new DatagramPacket(bytes, bytes.length, InetAddress.getByName(iparry[j]), 1000);
            ds.send(dp);
        }
        ds.close();
    }
}
class endserver implements Runnable {
    private final JFrame servermain;
    private final JFrame client;

    public endserver(JFrame servermain, JFrame client) {
        this.servermain = servermain;
        this.client = client;
    }

    @Override
    public void run() {
        servermain.setVisible(false);
        client.setVisible(true);
        servermain.dispose();
        client.dispose();
        System.exit(0);
    }
}