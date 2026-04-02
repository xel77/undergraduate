import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import javax.swing.SwingWorker;

public class Main {
    private JFrame frame;

    public static void main(String[] args) {
        Main mainInstance = new Main();
        mainInstance.createAndShowGUI();
    }

    public void createAndShowGUI() {
        frame = new JFrame("通讯");
        frame.setSize(320, 150);
        ImageIcon coin = new ImageIcon("微信图片_20240419142656.png");
        frame.setIconImage(coin.getImage());
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        JButton serverButton = new JButton("本机作为服务器");
        JButton clientButton = new JButton("本机作为客户端");
        frame.setLayout(new BorderLayout());
        JPanel panel = new JPanel();
        panel.add(serverButton);
        panel.add(clientButton);
        frame.add(panel, BorderLayout.BEFORE_FIRST_LINE);
        frame.setVisible(true);

        serverButton.addActionListener(e -> {
            serverButton.setEnabled(false);
            SwingWorker<Void, Void> worker = new SwingWorker<>() {
                @Override
                protected Void doInBackground() throws Exception {
                    frame.setVisible(false);
                    server server = new server(frame);
                    server.creatui();
                    return null;
                }

                @Override
                protected void done() {
                    serverButton.setEnabled(true);
                    reloadMain();
                }
            };
            worker.execute();
        });
        clientButton.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                frame.setVisible(false);
                client client = new client();
                client.creatui();
            }
        });
        clientButton.addActionListener(e -> {
            clientButton.setEnabled(false);
            SwingWorker<Void, Void> worker = new SwingWorker<>() {
                @Override
                protected Void doInBackground(){
                    frame.setVisible(false);
                    client client = new client();
                    client.creatui();
                    return null;
                }

                @Override
                protected void done() {
                    clientButton.setEnabled(true);
                    reloadMain();
                }
            };
            worker.execute();
        });
    }

    private void reloadMain() {
        SwingUtilities.invokeLater(() -> {
            frame.dispose();
            main(new String[]{});
        });
    }
}