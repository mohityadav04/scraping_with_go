go build -o main . && ./main

golint <file_name>

lower case named variable/func will not be imported to other packages



https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/


https://dev.to/dan_mcm_/lucas---a-webscraper-in-go-3jkc
https://github.com/gocolly/colly/tree/master/_examples
http://go-colly.org/


https://appliedgo.net/regexp/



Requirements:
1. Strict REST principle
2. Store in MongoDb
3. Dockerize using docker compose having db image

Assumptions:
1. reviews are ratings

Stages(break, sequences, research every, **incremental development(TDD if possible)**, fast):
1. Learn Go, Mux, Colly: Done
2. Scrape, print.
3. Choose document store, store it.
4. Dock
5. Edge cases, code refactoring(modularity, consistency), remove magic numbers, tests


TEST URLs:
1. https://www.amazon.com/GOOVI-Upgraded-Multiple-Self-Charging-Medium-Pile/dp/B081ZV69YR/ref=gbps_tit_m-9_475e_f68276a9?smid=A1JPKOLVI9WTHD&pf_rd_p=5d86def2-ec10-4364-9008-8fbccf30475e&pf_rd_s=merchandised-search-9&pf_rd_t=101&pf_rd_i=15529609011&pf_rd_m=ATVPDKIKX0DER&pf_rd_r=86GC220VSG1F45G1SHYD&spLa=ZW5jcnlwdGVkUXVhbGlmaWVyPUEzOU9RWEIzRUVPSFNVJmVuY3J5cHRlZElkPUEwOTQ5NDg2Mko0V0c4UUlZUVlIRiZlbmNyeXB0ZWRBZElkPUEwOTc5ODY1MzYyR1dONkZQU05FNSZ3aWRnZXROYW1lPXNwX2diX21haW5fc3VwcGxlJmFjdGlvbj1jbGlja1JlZGlyZWN0JmRvTm90TG9nQ2xpY2s9dHJ1ZQ==#customerReviews

2. https://www.amazon.com/dp/B01MSYY5X5/ref=sspa_dk_detail_0?psc=1&pd_rd_i=B01MSYY5X5&pd_rd_w=kZDyS&pf_rd_p=48d372c1-f7e1-4b8b-9d02-4bd86f5158c5&pd_rd_wg=L5Rnx&pf_rd_r=TQW2WNQ435HVQ3QAKJNC&pd_rd_r=98b411e2-20d0-4e65-82fe-1ec15670f5f7&spLa=ZW5jcnlwdGVkUXVhbGlmaWVyPUEyNFoyUjBGWTlDNlJQJmVuY3J5cHRlZElkPUEwNDQ1NzMzMkpCVzZUNVIwTzRJNiZlbmNyeXB0ZWRBZElkPUEwNzMxNzAzQUdSSU5WWEdRVDFMJndpZGdldE5hbWU9c3BfZGV0YWlsJmFjdGlvbj1jbGlja1JlZGlyZWN0JmRvTm90TG9nQ2xpY2s9dHJ1ZQ==

3. https://www.amazon.com/dp/B01MF9DLQW/ref=dp_cr_wdg_tit_nw_mr