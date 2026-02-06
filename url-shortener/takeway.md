# Takeways

## Different between Encoding, Hashing and Encryption.

### 1. Encoding

Encoding is the process of transformation one data format into another data format usign a scheme that is publicly available. It is used to improve the readability of data. It is easily reversible

#### Tpyes of Encoding

##### 1. Base64 Encoding:

It is used to encode binary data into ASCII string format.

**How it works?**
Base64 encode binary data by dividing the input data into 3-byte blocks and then encode each block into four 6-bit groups. Each 6-bit group is then mapped to a Base64 character.

**Character Set:**

- A-Z (26 characters)
- a-z (26 characters)
- 0-9 (10 characters)
- - and / (2 characters)

**Padding:**
If the input data is not a multiple of 3 bytes, then padding is added to make it a multiple of 3 bytes. Padding is done by adding one or two equal signs at the end of the encoded data.

**Why it is used?**
1- Text safe encoding: it ensure binary data such as images, audio, video, etc. can be transmitted as text over text-based protocols like HTTP, FTP, SMTP, etc.
2- Data Integrity: It is used to ensure data integrity during data transmission.

**Example:**

```
Input: "Hello"
binary: 01001000 01100101 01101100 01101100 01101111
divide into 3-byte blocks: 01001000 01100101 01101100 | 01101100 01101111
encode each block into four 6-bit groups: 010010 000110 010101 101100 011011 000110 1111
map each 6-bit group to a Base64 character: SGVsbG8==
```

**Another encoding schemes:**

- Base62 Encoding: It is used to encode binary data into ASCII string format using 62 characters, removes / and + characters to avoid confusion.
- Base58 Encoding: It is used to encode binary data into ASCII string format using 58 characters, removes O, 0, I, l characters to avoid confusion.

The alogrithm for Base62 and Base58 encoding is similar to Base64 encoding.

example: Hello

1- We convert the binary data into a number: 01001000 01100101 01101100 01101100 01101111 = 72 101 108 108 111
2- We accumulate the number: 72 _ 256^4 + 101 _ 256^3 + 108 _ 256^2 + 108 _ 256 + 111 = 1194633975
3- We take the number and convert it to base62 or base58 by dividing the number by 62 or 58 and taking the remainder.
e.g 1194633975 % 62 = 43, 1194633975 / 62 = 19284903, 19284903 % 62 = 3, 19284903 / 62 = 311370, 311370 % 62 = 30, 311370 / 62 = 5020, 5020 % 62 = 20, 5020 / 62 = 81, 81 % 62 = 19, 81 / 62 = 1, 1 % 62 = 1
output: 1Z3a20

### 2. Hashing

Hashing is the process of transforming input data into a fixed-size string of bytes using a hash function, usually represent the dat in a way that is suitable for integirty or data lookup. It is not reversible.

#### Properties of Hashing

1. Deterministic: The same input will always produce the same output.
2. Fast: It should be fast to compute the hash value.
3. Fixed-size output: The output should be of fixed size, regardless of the input size.
4. Irreversible: It should be computationally infeasible to reverse the hash value to get the original input.
5. Collision-resistant: It should be computationally infeasible to find two different inputs that produce the same hash value.

#### Use Cases of Hashing

1. Password Storage: Hashing is used to store passwords securely. Instead of storing the actual password, the hash value of the password is stored in the database. When a user logs in, the hash value of the entered password is compared with the stored hash value.

2. Data Integrity: Hashing is used to ensure data integrity during data transmission. The sender computes the hash value of the data and sends it along with the data. The receiver computes the hash value of the received data and compares it with the received hash value to check if the data has been tampered with.

3. Digital Signatures: Hashing is used in digital signatures to ensure the authenticity and integrity of the message. The sender computes the hash value of the message and encrypts it with their private key to create a digital signature. The receiver decrypts the digital signature with the sender's public key and compares the hash value of the message with the decrypted hash value to verify the authenticity and integrity of the message.

#### Hashing Algorithms

##### 1. Cryptographic Hash Functions

Cryptographic hash functions are designed to be secure, meaning they are resistant to attacks and collisions. They are used in security-critical applications like password storage, digital signatures, and data integrity verification.

1. MD5 (Message Digest Algorithm 5): It produces a 128-bit hash value and is no longer considered secure due to vulnerabilities in the algorithm. MD5 consider faster than SHA-1 and it's still used in some non-security-critical applications like checksums and data integrity verification.

2. SHA-1 (Secure Hash Algorithm 1): It produces a 160-bit hash value and is no longer considered secure due to vulnerabilities in the algorithm. SHA-1 is still used in some legacy systems and applications.

3. SHA-256 (Secure Hash Algorithm 256): It produces a 256-bit hash value and is considered secure for most applications. It is widely used in blockchain technology, digital signatures, and data integrity verification.

4. SHA-3 (Secure Hash Algorithm 3): It produces hash values of various sizes (224, 256, 384, 512 bits) and is considered secure. It is the latest member of the Secure Hash Algorithm family.

##### 2. Non-Cryptographic Hash Functions

Non-cryptographic hash functions are designed for speed and efficiency, rather than security. They are used in applications like hash tables, data lookup, and checksums.

1. MurmurHash: It is a non-cryptographic hash function that is fast and efficient. It is widely used in hash tables and data lookup applications.

2. CityHash: It is a non-cryptographic hash function that is optimized for speed and efficiency. It is used in applications like hash tables and data lookup.

3. FNV (Fowler-Noll-Vo): It is a non-cryptographic hash function that is simple and fast. It is used in applications like hash tables and checksums.

##### 3. Checksum Hash Functions

Checksum hash functions are used to detect errors in data transmission. They are not designed for security but for error detection.

1. CRC32 (Cyclic Redundancy Check 32): It produces a 32-bit hash value and is used for error detection in data transmission.

2. Adler-32: It produces a 32-bit hash value and is used for error detection in data transmission.

3. Fletcher's Checksum: It produces a 16-bit or 32-bit hash value and is used for error detection in data transmission.

### 3. Encryption

Encryption is used to transform data into a format that is unreadable format to prevent unauthorized access. It is reversible, meaning the encrypted data can be decrypted back to its original form using a decryption key.

#### Key Concepts of Encryption

1. Plaintext: The original data that is to be encrypted.
2. Ciphertext: The encrypted data that is produced after encryption.
3. Encryption Algorithm (Cipher suite): The algorithm used to encrypt the data.
4. Decryption Algorithm: The algorithm used to decrypt the data.
5. Encryption Key: The key used to encrypt the data.
6. Decryption Key: The key used to turn ciphertext back into plaintext. Depending on the encryption method, this can be the same as the encryption key (symmetric encryption) or different (asymmetric encryption).

#### Types of Encryption

##### 1. Symmetric Encryption:

In symmetric encryption, the same key is used for both encryption and decryption. The key is shared between the sender and the receiver.

**Advantages:**

1. Fast: Symmetric encryption is faster than asymmetric encryption.
2. Simple: Symmetric encryption is simpler than asymmetric encryption.

**Disadvantages:**

1. Key Distribution: The key must be securely shared between the sender and the receiver.
2. Key Management: The key must be managed securely to prevent unauthorized access.

**Examples:**

- AES (Advanced Encryption Standard): It is a symmetric encryption algorithm that is widely used for data encryption. It supports key sizes of 128, 192, and 256 bits.
- DES (Data Encryption Standard): It is a symmetric encryption algorithm that is no longer considered secure due to its small key size (56 bits).
- 3DES (Triple Data Encryption Standard): It is a symmetric encryption algorithm that applies DES three times to increase the key size to 168 bits.
- chacha20: It is a symmetric encryption algorithm that is designed to be secure and efficient (used in PASETO an alternative to JWT).

##### 2. Asymmetric Encryption:

In asymmetric encryption, two keys are used: a public key for encryption and a private key for decryption. The public key is shared with everyone, while the private key is kept secret.

**Advantages:**

1. Key Distribution: The public key can be shared with everyone, while the private key is kept secret.
2. Key Management: The private key is kept secret and is not shared with anyone.

**Disadvantages:**

1. Slow: Asymmetric encryption is slower than symmetric encryption.
2. Complexity: Asymmetric encryption is more complex than symmetric encryption.

**Examples:**

- RSA (Rivest-Shamir-Adleman): It is an asymmetric encryption algorithm that is widely used for secure communication, digital signatures, and data encryption.
- ECC (Elliptic Curve Cryptography): It is an asymmetric encryption algorithm that is based on the algebraic structure of elliptic curves. It is used in applications where key size and speed are important.
- Diffie-Hellman Key Exchange: It is an asymmetric encryption algorithm that is used to securely exchange keys over an insecure channel.
