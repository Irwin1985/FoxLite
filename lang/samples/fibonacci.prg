FUNCTION FIBO(N)
	IF N <= 1
		RETURN N
	ENDIF
	RETURN FIBO(N-1) + FIBO(N-2)
ENDFUNC
