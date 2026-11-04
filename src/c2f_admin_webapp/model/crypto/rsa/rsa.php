<?php
$plaintext = 'Hello PHP';
 
 
$publicKey = openssl_pkey_get_public(file_get_contents('/Users/nirbhay/KeyStore/BillingEnginePGP/public_key.cer'));
if(!$publicKey){
	echo openssl_error_string();
	die("");
}
$a_key = openssl_pkey_get_details($publicKey);
if(!$a_key){
	echo openssl_error_string();
	die("");
}

$chunkSize = ceil($a_key['bits'] / 8) - 11;
$output = '';
 
while ($plaintext)
{
    $chunk = substr($plaintext, 0, $chunkSize);
    $plaintext = substr($plaintext, $chunkSize);
    $encrypted = '';
    if (!openssl_public_encrypt($chunk, $encrypted, $publicKey))
    {
        die('Failed to encrypt data');
    }
    $output .= $encrypted;
}
openssl_free_key($publicKey);
 
$encrypted = base64_encode($output);

print_r($encrypted);
echo "\n";
?>