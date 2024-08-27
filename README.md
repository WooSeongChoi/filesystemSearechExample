# filesystemSearchExample
## 개요
- 간단하게 기준점의 자식 directory를 대상으로 goroutine으로 Walk하는 것을 검증

## 한계점
- 자식 directory 하위 계산량이 균등하지 않으면 개선의 폭이 좁아진다.

## 개선이 필요한 점
- 자식 directory 조회시 goroutine 사용할 것.