# 1. Generate a bianry file according to your meta file:
## 1.1 A meta file example:
```
{
  "FieldInfos" : [
    {"FieldLen": 8, "FieldName": "TimeStamp", "FieldType": "int64", "InitVal": "1"},
    {"FieldLen": 8, "FieldName": "OrderEntryTime", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 8, "FieldName": "OrderExpiryTime", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 8, "FieldName": "OrderNumber", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 8, "FieldName": "OrderPrice", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 8, "FieldName": "OrderExecutedQuantity", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 8, "FieldName": "OrderQuantity", "FieldType": "int64", "InitVal": "1"},
	{"FieldLen": 2, "FieldName": "isix", "FieldType": "uint16", "InitVal": "1"},
	{"FieldLen": 5, "FieldName": "PBUID", "FieldType": "string", "InitVal": "12345"},
	{"FieldLen": 10, "FieldName": "InvestorAccountID", "FieldType": "string", "InitVal": "12345789"},
	{"FieldLen": 1, "FieldName": "AccountTypeCode", "FieldType": "string","InitVal": "C"},
	{"FieldLen": 5, "FieldName": "BranchID", "FieldType": "string", "InitVal": "ABCDE"},
	{"FieldLen": 16, "FieldName": "PbuInternalOrderNum", "FieldType": "string", "InitVal": "hjiangheng"},
	{"FieldLen": 12, "FieldName": "Text", "FieldType": "string", "InitVal": "hello world"},
	{"FieldLen": 5, "FieldName": "ClearingParticipantID", "FieldType": "string", "InitVal": "hello"},
	{"FieldLen": 1, "FieldName": "AuditTypeCode", "FieldType": "string", "InitVal": "B"},
	{"FieldLen": 2, "FieldName": "OrderTypeCode", "FieldType": "string", "InitVal": "HA"}, 
	{"FieldLen": 2, "FieldName": "TrdRestrictionTypeCode", "FieldType": "string", "InitVal": "JH"},
	{"FieldLen": 1, "FieldName": "OrderStatus", "FieldType": "string", "InitVal": "D"},
	{"FieldLen": 1, "FieldName": "BuySellCoder", "FieldType": "string", "InitVal": "H"},
	{"FieldLen": 1, "FieldName": "OrderRestrictionCode", "FieldType": "string", "InitVal": "G"},
	{"FieldLen": 1, "FieldName": "ShortSellCheckingFlag", "FieldType": "string", "InitVal": "U"},
	{"FieldLen": 1, "FieldName": "TransactionMaintainCode", "FieldType": "string", "InitVal": "I"},
	{"FieldLen": 1, "FieldName": "CreditTag", "FieldType": "string", "InitVal": "O"},
	{"FieldLen": 5, "FieldName": "Filler", "FieldType": "string", "InitVal": "FILLE"} ]
}
```
the meaning of each filed:
+ FieldLen: The length of this field, 8 stand for 8 bytes;
+ FieldName: The name of this field, tell the meaning of this field;
+ FieldType: The data type of this field;
+ InitVal: The value of this field when write to a bianry file;

## 1.2 Generate the bianry file
When you define out the meta file, you can generate the bianry file:
```
Usage:
  binarygen encode [flags]

Flags:
  -h, --help          help for encode
  -m, --meta string   meta data file path (default "./configs/meta.json")
  -n, --num int       number of records in the output file (default 10)
  -o, --out string    output data file path (default "./out.log")

Global Flags:
      --config string   config file (default is $HOME/.binarygen.yaml)
```

```
./binarygen encode -m yourmetafilepath -n 100 -o ./out.log
```

# 2. Decode the binary file
Use the same meta data, you can decode the bianry file.
```
Usage:
  binarygen decode [flags]

Flags:
  -h, --help          help for decode
  -i, --in string     binary data file path (default "./in.log")
  -m, --meta string   meta data file path (default "./configs/meta.json")
  -o, --out string    output data file path (default "./out.log")

Global Flags:
      --config string   config file (default is $HOME/.binarygen.yaml)
```

```
./binarygen decode -m yourmetafilepath -i ./out.log -o ./bout.txt
```
