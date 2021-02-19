from src.book.fox_lite_token import TokenType
from src.book.fox_lite_ast import (
    Program,
    Block,
    BinaryOp,
    UnaryOp,
    FunctionDecl,
    FunctionCall,
    DoWhile,
    IfStatement,
    Identifier,
    Boolean,
    Integer,
    String,
    Null,
    VariableDecl,
    Assignment,
    ReturnStmt,
    PrintStmt,
)

"""
El parser se encarga de validar la gramática y crear el AST. La gramática completa es la siguiente:
program ::= (statement)+ EOF
block ::= (statement)+ EOF
statementList ::= statement (LBREAK statement)*
statement ::= variableDeclarationStmt | functionDeclarationStmt | doWhileStmt | ifStmt
variableDeclarationStmt ::= publicDeclaration | localDeclaration | privateDeclaration
publicDeclaration ::= 'public' identifier LBREAK
localDeclaration ::= 'local' identifier LBREAK
privateDeclaration ::= 'private' identifier LBREAK
functionDeclaration ::= 'function' '(' ( arguments )? ')' LBREAK block 'endfunc'
doWhileStmt ::= 'do while' expression LBREAK block 'enddo'
ifStmt ::= 'if' expression LBREAK block ('else' block)? 'endif'
assignment ::= identifier '=' expression LBREAK
expression ::= logicOr
logicOr ::= logicAnd ('or' logicAnd )*
logicAnd ::= equality ('and' equality )*
equality ::= comparison ( ('==' | '!=') comparison )*
comparison ::= term ( ( '<' | '>' | '<=' | '>=' ) term )*
term ::= factor ( ( '+' | '-' ) factor )*
factor ::= unary ( ( '*' | '/' ) unary )*
unary ::= primary ( ( '!' | '-') primary )*
primary ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER
"""


class Parser:
    def __init__(self, lexer):
        self.lexer = lexer
        self.cur_token = None
        self.peek_token = None
        self.errors = []
        self.next_token()  # Avanza el primer token
        self.next_token()  # Avanza el segundo token

    def next_token(self):
        self.cur_token = self.peek_token
        self.peek_token = self.lexer.next_token()

    def eat(self, token_type):
        if self.cur_token.type == token_type:
            self.next_token()
        else:
            msg = f"Se esperaba el token: '{token_type}', y se obtuvo '{self.cur_token.type} en su lugar.'"
            self.errors.append(msg)

    """
    program ::= (statement)+ EOF
    """
    def program(self):
        program = Program()
        while self.cur_token.type != TokenType.EOF:
            statement = self.statement()
            if statement is not None:
                program.statements.append(statement)

        return program

    """
    block ::= (statement)+ 'enddo' | 'endif' | 'endfunc'
    """
    def block(self, ending_block_token_type):
        block = Block()

        if self.cur_token.type == TokenType.LBREAK:
            self.eat(TokenType.LBREAK)

        while self.cur_token.type != ending_block_token_type:
            statement = self.statement()
            if statement is not None:
                block.statements.append(statement)

        self.eat(ending_block_token_type)

        if self.cur_token.type == TokenType.LBREAK:
            self.eat(TokenType.LBREAK)

        return block

    def statement(self):
        if self.cur_token.type == TokenType.PUBLIC:
            return self.public_decl()
        elif self.cur_token.type == TokenType.LOCAL:
            return self.local_decl()
        elif self.cur_token.type == TokenType.PRIVATE:
            return self.private_decl()
        elif self.cur_token.type == TokenType.FUNCTION:
            return self.function_decl()
        elif self.cur_token.type == TokenType.DO:
            return self.do_while_decl()
        elif self.cur_token.type == TokenType.IF:
            return self.if_statement()
        elif self.cur_token.type == TokenType.IDENT and self.peek_token.type == TokenType.ASSIGN:
            return self.assignment()
        elif self.cur_token.type == TokenType.RETURN:
            return self.return_statement()
        elif self.cur_token.type == TokenType.PRINT:
            return self.print_statement()
        else:
            return self.expression()
    """
        public_decl = 'public' identifier
    """
    def public_decl(self):
        self.eat(TokenType.PUBLIC)
        var_token = self.identifier()
        return VariableDecl(token=var_token, scope='public')
    """
        local_decl = 'local' identifier
    """
    def local_decl(self):
        self.eat(TokenType.LOCAL)
        var_token = self.identifier()
        return VariableDecl(token=var_token, scope='local')
    """
        private = 'private' identifier
    """
    def private_decl(self):
        self.eat(TokenType.PRIVATE)
        var_token = self.identifier()
        return VariableDecl(token=var_token, scope='private')
    """
        function_decl ::= 'function' identifier '(' ( parameters )? ')'
    """
    def function_decl(self):
        self.eat(TokenType.FUNCTION)
        func = FunctionDecl(name=self.identifier())

        if self.cur_token.type == TokenType.LPAREN:
            self.eat(TokenType.LPAREN)

            if self.cur_token.type != TokenType.RPAREN:
                func.params = self.parse_parameters()

            self.eat(TokenType.RPAREN)

        func.body = self.block(ending_block_token_type=TokenType.ENDFUNC)

        return func
    """
        do_while_decl ::= 'do while' expression block 'enddo'
    """
    def do_while_decl(self):
        self.eat(TokenType.DO)
        self.eat(TokenType.WHILE)
        do_while = DoWhile()
        do_while.condition = self.expression()
        do_while.block = self.block(ending_block_token_type=TokenType.ENDDO)

        return do_while
    """
       if_statement ::= 'if' condition block ( 'else' block )? 'endif' 
    """
    def if_statement(self):
        self.eat(TokenType.IF)
        if_stmt = IfStatement()
        if_stmt.condition = self.expression()
        if_stmt.consequence = Block()

        while self.cur_token.type not in (TokenType.ELSE, TokenType.ENDIF):
            statement = self.statement()
            if statement is not None:
                if_stmt.consequence.statements.append(statement)

        if self.cur_token.type == TokenType.ELSE:
            self.eat(TokenType.ELSE)
            if_stmt.alternative = self.block(ending_block_token_type=TokenType.ENDIF)
        else:
            self.eat(TokenType.ENDIF)

        return if_stmt
    """
        assignment ::= identifier '=' expression LBREAK
    """
    def assignment(self):
        ident = self.identifier()
        self.eat(TokenType.ASSIGN)
        value = self.expression()
        return Assignment(token=ident, value=value)

    """
        return_statement ::= 'return' (expression)?
    """
    def return_statement(self):
        return_stmt = ReturnStmt()
        self.eat(TokenType.RETURN)

        if self.cur_token.type != TokenType.LBREAK:
            return_stmt.value = self.expression()
        else:
            return_stmt.value = Boolean(value=True)

        return return_stmt

    def print_statement(self):
        print_stmt = PrintStmt()
        self.eat(TokenType.PRINT)
        print_stmt.arguments.append(self.expression())

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            print_stmt.arguments.append(self.expression())

        return print_stmt

    """
    expression ::= logicOr
    logicOr    ::= logicAnd ('or' logicAnd )*
    logicAnd   ::= equality ('and' equality )*
    equality   ::= comparison ( ('==' | '!=') comparison )*
    comparison ::= term ( ( '<' | '>' | '<=' | '>=' ) term )*
    term       ::= factor ( ( '+' | '-' ) factor )*
    factor     ::= unary ( ( '*' | '/' ) unary )*
    unary      ::= primary ( ( '!' | '-') primary )*
    primary    ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER    
    """
    def expression(self):
        exp = self.logic_or()

        if self.cur_token.type == TokenType.LBREAK:
            self.eat(TokenType.LBREAK)

        return exp
    """
    logicOr ::= logicAnd ('or' logicAnd )* 
    """
    def logic_or(self):
        or_exp = self.logic_and()
        while self.cur_token.type == TokenType.OR:
            tok = self.cur_token
            self.eat(TokenType.OR)
            or_exp = BinaryOp(left=or_exp, operator=tok.value, right=self.logic_and())

        return or_exp
    """
    logicAnd ::= equality ('and' equality )*  
    """
    def logic_and(self):
        and_exp = self.equality()
        while self.cur_token.type == TokenType.AND:
            tok = self.cur_token
            self.eat(TokenType.AND)
            and_exp = BinaryOp(left=and_exp, operator=tok.value, right=self.equality())

        return and_exp
    """
    equality ::= comparison ( ('==' | '!=') comparison )*  
    """
    def equality(self):
        comp_exp = self.comparison()
        while self.cur_token.type in (TokenType.EQUAL, TokenType.NOT_EQUAL):
            tok = self.cur_token
            self.eat(tok.type)  # '==' ó '!='
            comp_exp = BinaryOp(left=comp_exp, operator=tok.value, right=self.comparison())

        return comp_exp
    """
    comparison ::= term ( ( '<' | '>' | '<=' | '>=' ) term )*  
    """
    def comparison(self):
        term_exp = self.term()
        while self.cur_token.type in (TokenType.LESS, TokenType.LESS_EQ,
                                      TokenType.GREATER, TokenType.GREATER_EQ):
            tok = self.cur_token
            self.eat(tok.type)  # '<' ó '>' ó '<=' ó '>='
            term_exp = BinaryOp(left=term_exp, operator=tok.value, right=self.term())

        return term_exp
    """
    term ::= factor ( ( '+' | '-' ) factor )*   
    """
    def term(self):
        factor_exp = self.factor()
        while self.cur_token.type in (TokenType.PLUS, TokenType.MINUS):
            tok = self.cur_token
            self.eat(tok.type)  # '+' ó '-'
            factor_exp = BinaryOp(left=factor_exp, operator=tok.value, right=self.factor())

        return factor_exp
    """
    factor ::= unary ( ( '*' | '/' ) unary )*
    """
    def factor(self):
        unary_exp = self.unary()
        while self.cur_token.type in (TokenType.MUL, TokenType.DIV):
            tok = self.cur_token
            self.eat(tok.type)  # '*' ó '/'
            unary_exp = BinaryOp(left=unary_exp, operator=tok.value, right=self.unary())

        return unary_exp
    """
    unary ::= primary ( ( '!' | '-') primary )*  
    """
    def unary(self):
        if self.cur_token.type in (TokenType.NOT, TokenType.MINUS):
            operator = self.cur_token.value
            self.eat(self.cur_token.type)  # '!' ó '-'
            return UnaryOp(operator=operator, right=self.unary())
        else:
            return self.call()
    """
    call ::= function_call | primary    
    """
    def call(self):
        primary = self.primary()
        if self.cur_token.type == TokenType.LPAREN:
            function_name = primary
            primary = FunctionCall()
            primary.name = function_name

            self.eat(TokenType.LPAREN)

            if self.cur_token.type != TokenType.RPAREN:
                primary.arguments = self.parse_arguments()

            self.eat(TokenType.RPAREN)

        return primary
    """
    primary ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER | expression  
    """
    def primary(self):
        tok = self.cur_token

        if tok.type in (TokenType.TRUE, TokenType.FALSE):
            self.eat(tok.type)  # .T. ó .F.
            return Boolean(value=(tok.type == TokenType.TRUE))

        elif tok.type == TokenType.NULL:
            self.eat(TokenType.NULL)
            return Null()

        elif tok.type == TokenType.INT:
            self.eat(TokenType.INT)
            return Integer(value=int(tok.value))

        elif tok.type == TokenType.STRING:
            self.eat(TokenType.STRING)
            return String(value=tok.value)

        elif tok.type == TokenType.IDENT:
            return self.identifier()

        elif tok.type == TokenType.LPAREN:
            self.eat(TokenType.LPAREN)
            exp = self.expression()
            self.eat(TokenType.RPAREN)
            return exp

    def parse_arguments(self):
        args = [self.expression()]

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            args.append(self.expression())

        return args

    def parse_parameters(self):
        identifiers = [self.identifier()]

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            identifiers.append(self.identifier())

        return identifiers

    def identifier(self):
        name = self.cur_token.value
        self.eat(TokenType.IDENT)
        return Identifier(value=name)

    def parse(self):
        program = self.program()
        if self.cur_token.type != TokenType.EOF:
            self.errors.append(f'Faltaron tokens por analizar')

        return program
