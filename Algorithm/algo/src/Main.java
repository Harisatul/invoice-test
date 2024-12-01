import java.util.ArrayList;
import java.util.List;

public class Main {
    public static void main(String[] args) {
        int I = 3;
        int t = 6;
        List<List<Integer>> hasil = cariKombinasi(I, t);
        System.out.println(hasil);
    }

    private static void cari(int mulai, int panjang, int target, List<Integer> kombinasi, List<List<Integer>> hasil) {
        if (panjang == 0 && target == 0) {
            hasil.add(new ArrayList<>(kombinasi));
            return;
        }

        if (panjang == 0 || target < 0) {
            return;
        }

        for (int i = mulai; i <= 9; i++) {
            kombinasi.add(i);
            cari(i + 1, panjang - 1, target - i, kombinasi, hasil);
            kombinasi.remove(kombinasi.size() - 1);
        }
    }

    public static List<List<Integer>> cariKombinasi(int I, int t) {
        List<List<Integer>> hasil = new ArrayList<>();
        cari(1, I, t, new ArrayList<>(), hasil);
        return hasil;
    }
}