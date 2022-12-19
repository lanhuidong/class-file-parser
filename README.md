# class-file-parser
Java class文件解析器

### Build
```shell go run main.go -file fileName.class```

### Class文件格式
| 类型 | 名称 | 数量 |
|:---|:---|:---|
| u4 | magic | 1 |
| u2 | minor_version | 1 |
| u2 | major_version | 1 |
| u2 | constant_pool_count | 1 |
| cp_info | constant_pool | constant_pool_count - 1 |
| u2 | access_flags | 1 |
| u2 | this_class | 1 |
| u2 | super_class | 1 |
| u2 | interfaces_count | 1 |
| u2 | interfaces | interfaces_count |
| u2 | fields_count | 1 |
| field_info | fields | fields_count |
| u2 | methods_count | 1 |
| method_info | methods | methods_count |
| u2 | attributes_count | 1 |
| attribute_info | attributes | attributes_count |

### 常量池中数据类型结构总表
| 常量 | 项目 | 类型 | 描述 |
|:---|:---|:---|:---|
| CONSTANT_Utf8_info | tag | u1 | 值为1 |
| CONSTANT_Utf8_info | length | u2 | UTF-8字符串的字节数 |
| CONSTANT_Utf8_info | bytes | length个u1 | UTF-8字符串 |
| CONSTANT_Integer_info | tag | u1 | 值为3 |
| CONSTANT_Integer_info | bytes | u4 | 大端存储的int值 |
| CONSTANT_Float_info | tag | u1 | 值为4 |
| CONSTANT_Float_info | bytes | u4 | 大端存储的float值 |
| CONSTANT_Long_info | tag | u1 | 值为5 |
| CONSTANT_Long_info | bytes | u4 | 大端存储的long值 |
| CONSTANT_Double_info | tag | u1 | 值为6 |
| CONSTANT_Double_info | bytes | u4 | 大端存储的double值 |
| CONSTANT_Class_info | tag | u1 | 值为7 |
| CONSTANT_Class_info | index | u2 | 指向全限定名常量项CONSTANT_Utf8_info的索引 |
| CONSTANT_String_info | tag | u1 | 值为8 |
| CONSTANT_String_info | index | u2 | 指向字符串字面量CONSTANT_Utf8_info的索引 |
| CONSTANT_Fieldref_info | tag | u1 | 值为9 |
| CONSTANT_Fieldref_info | index | u2 | 指向声明字段的类或者接口描述符CONSTANT_Class_info的索引项 |
| CONSTANT_Fieldref_info | index | u2 | 指向名称及类型描述符CONSTANT_NameAndType_info的索引项 |
| CONSTANT_Methodref_info | tag | u1 | 值为10 |
| CONSTANT_Methodref_info | index | u2 | 指向声明方法的类描述符CONSTANT_Class_info的索引项 |
| CONSTANT_Methodref_info | index | u2 | 指向名称及类型描述符CONSTANT_NameAndType_info的索引项 |
| CONSTANT_InterfaceMethodref_info | tag | u1 | 值为11 |
| CONSTANT_InterfaceMethodref_info | index | u2 | 指向声明方法的接口描述符CONSTANT_Class_info的索引项 |
| CONSTANT_InterfaceMethodref_info | index | u2 | 指向名称及类型描述符CONSTANT_NameAndType_info的索引项 |
| CONSTANT_NameAndType_info | tag | u1 | 值为12 |
| CONSTANT_NameAndType_info | index | u2 | 指向该字段或方法名称常量项的索引 |
| CONSTANT_NameAndType_info | index | u2 | 指向该字段或方法描述符常量项的索引 |
| CONSTANT_MethodHandle_info | tag | u1 | 值为15 |
| CONSTANT_MethodHandle_info | reference_kind | u1 | 取值区间[1, 9]，它决定了方法句柄的类型。方法句柄类型值表示方法句柄的字节码行为 |
| CONSTANT_MethodHandle_info | reference_index | u2 | 值必须是对常量池的有效索引 |
| CONSTANT_MethodType_info | tag | u1 | 值为16 |
| CONSTANT_MethodType_info | descriptor_index | u2 | 值必须是对CONSTANT_Utf8_info的有效索引，表示方法的描述符 |
| CONSTANT_Dynamic_info | tag | u1 | 值为17 |
| CONSTANT_Dynamic_info | bootstrap_method_attr_index | u2 | 值必须是对当前Class文件中引导方法表bootstrap_methods[]数据的有效索引 |
| CONSTANT_Dynamic_info | name_and_type_index | u2 | 值必须是对CONSTANT_NameAndType_info的有效索引，表示方法名和方法描述符 |
| CONSTANT_InvokeDynamic_info | tag | u1 | 值为18 |
| CONSTANT_InvokeDynamic_info | bootstrap_method_attr_index | u2 | 值必须是对当前Class文件中引导方法表bootstrap_methods[]数据的有效索引 |
| CONSTANT_InvokeDynamic_info | name_and_type_index | u2 | 值必须是对CONSTANT_NameAndType_info的有效索引，表示方法名和方法描述符 |
| CONSTANT_Module_info | tag | u1 | 值为19 |
| CONSTANT_Module_info | name_index | u2 | 值必须是对CONSTANT_Utf8_info的有效索引，表示模块名称 |
| CONSTANT_Package_info | tag | u1 | 值为20 |
| CONSTANT_Package_info | name_index | u2 | 值必须是对CONSTANT_Utf8_info的有效索引，表示包名 |

注：CONSTANT_Long_info和CONSTANT_Double_info在常量池中占2个长度