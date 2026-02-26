#!/bin/bash

echo "=== 检查 admin 角色权限 ==="
mysql -h192.168.230.131 -uroot -p'Leo@123456' watchalert1 -e "SELECT id, name, permissions FROM user_roles WHERE name = 'admin' LIMIT 1;" 2>/dev/null

echo ""
echo "=== 检查知识库相关权限 ==="
mysql -h192.168.230.131 -uroot -p'Leo@123456' watchalert1 -e "SELECT permissions FROM user_roles WHERE name = 'admin' AND permissions LIKE '%knowledge%';" 2>/dev/null

echo ""
echo "=== 检查智能派单相关权限 ==="
mysql -h192.168.230.131 -uroot -p'Leo@123456' watchalert1 -e "SELECT permissions FROM user_roles WHERE name = 'admin' AND permissions LIKE '%assignment%';" 2>/dev/null
