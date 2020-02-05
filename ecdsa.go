package main

//func main() {
//	//创建曲线
//	curve := elliptic.P256()
//	//生成私匙
//	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//生成公钥
//	publicKey := privateKey.PublicKey
//	//对数据进行哈希运算
//	data := "666666666"
//	hash := sha256.Sum256([]byte(data))
//	//数据签名
//	//func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error) {
//	r, s, er := ecdsa.Sign(rand.Reader, privateKey, hash[:])
//	if er != nil {
//		fmt.Println(err)
//	}
//	//把r、s进行序列化传输
//	//1.传输
//	signature := append(r.Bytes(), s.Bytes()...)
//	//2.获取、定义两个辅助的BIG.INT
//	r1 := big.Int{}
//	s1 := big.Int{}
//	//3.拆分并赋值
//	r1.SetBytes(signature[:len(signature)/2])
//	s1.SetBytes(signature[len(signature)/2:])
//
//	//数据校验
//	//func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
//	status := ecdsa.Verify(&publicKey, hash[:], &r1, &s1)
//	fmt.Println(status)
//
//}
