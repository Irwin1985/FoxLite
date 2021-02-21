from src.book.fox_lite_token import TokenType
import src.book.fox_lite_ast as ast

"""
El parser se encarga de validar y crear el AST usando la siguiente gramática:
program         ::= statement (LBREAK statement)* EOF
block           ::= (LBREAK)? (statement)+ 'enddo' | 'endif' | 'endfunc' (LBREAK)?
statement       ::= public_decl
                | local_decl
                | private_decl 
                | function_decl 
                | do_while_decl
                | if_decl
                | assignment_decl
                | return_decl
                | print_decl
                | expression_decl
public_decl     ::= 'public' identifier LBREAK
local_decl      ::= 'local' identifier LBREAK
private_decl    ::= 'private' identifier LBREAK
function_decl   ::= 'function' '(' ( arguments )? ')' block 'endfunc'
do_while_decl   ::= 'do while' expression_decl block 'enddo'
if_decl         ::= 'if' expression_decl block ('else' block)? 'endif'
assignment_decl ::= identifier '=' expression_decl LBREAK
return_decl     ::= 'return' expression_decl
print_decl      ::= '?' expression_decl ( ',' expression_decl)*
expression_decl ::= logic_or
logic_or         ::= logic_and ('or' logic_and )*
logic_and        ::= equality ('and' equality )*
equality        ::= comparison ( ('==' | '!=') comparison )*
comparison      ::= term ( ( '<' | '>' | '<=' | '>=' ) term )*
term            ::= factor ( ( '+' | '-' ) factor )*
factor          ::= unary ( ( '*' | '/' ) unary )*
unary           ::= ( '!' | '-') unary | functionCall
functionCall    ::= primary ( '(' arguments? ')' )?
primary         ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER | '(' expression_decl ')'
"""


class Parser:
    """
    Parser: se encarga de validar la gramática y generar el árbol de sintaxis abstracta (AST).
    """
    def __init__(self, lexer):
        self.lexer = lexer
        self.cur_token = None
        self.peek_token = None
        self.errors = []
        self.next_token()  # Avanza el primer token
        self.next_token()  # Avanza el segundo token

    def next_token(self):
        """
        Avanza tanto el token actual como el siguiente token.
        :return:
        """
        self.cur_token = self.peek_token
        self.peek_token = self.lexer.next_token()

    def eat(self, token_type):
        """
        Verifica si el token dado es el mismo que el token actual para avanzarlos.
        :param token_type:
        :return: None
        """
        if self.cur_token.type == token_type:
            self.next_token()
        else:
            msg = f"Se esperaba el token: '{token_type}', y se obtuvo '{self.cur_token.type}' en su lugar."
            self.errors.append(msg)

    """
    program ::= statement (LBREAK statement)* EOF
    """
    def program(self):
        program = ast.Program()
        while self.cur_token.type != TokenType.EOF:
            statement = self.statement()
            if statement is not None:
                program.statements.append(statement)

        return program

    """
    block ::= (LBREAK)? (statement)+ 'enddo' | 'endif' | 'endfunc' (LBREAK)?
    """
    def block(self, ending_block_token_type):
        block = ast.Block()

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
    """
    statement   ::= public_decl
                | local_decl
                | private_decl 
                | function_decl 
                | do_while_decl
                | if_decl
                | assignment_decl
                | return_decl
                | print_decl
                | expression_decl
    """
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
            return self.if_decl()
        elif self.cur_token.type == TokenType.IDENT and self.peek_token.type == TokenType.ASSIGN:
            return self.assignment_decl()
        elif self.cur_token.type == TokenType.RETURN:
            return self.return_decl()
        elif self.cur_token.type == TokenType.PRINT:
            return self.print_decl()
        else:
            return self.expression_decl()
    """
        public_decl = 'public' identifier
    """
    def public_decl(self):
        self.eat(TokenType.PUBLIC)
        var_token = self.identifier()
        return ast.VariableDecl(token=var_token, scope='public')
    """
        local_decl = 'local' identifier
    """
    def local_decl(self):
        self.eat(TokenType.LOCAL)
        var_token = self.identifier()
        return ast.VariableDecl(token=var_token, scope='local')
    """
        private_decl = 'private' identifier
    """
    def private_decl(self):
        self.eat(TokenType.PRIVATE)
        var_token = self.identifier()
        return ast.VariableDecl(token=var_token, scope='private')
    """
        function_decl ::= 'function' identifier '(' ( parameters )? ')'
    """
    def function_decl(self):
        self.eat(TokenType.FUNCTION)
        func = ast.FunctionDecl(name=self.identifier())

        if self.cur_token.type == TokenType.LPAREN:
            self.eat(TokenType.LPAREN)

            if self.cur_token.type != TokenType.RPAREN:
                func.params = self.parse_parameters()

            self.eat(TokenType.RPAREN)

        func.body = self.block(ending_block_token_type=TokenType.ENDFUNC)

        return func
    """
        do_while_decl ::= 'do while' expression_decl block 'enddo'
    """
    def do_while_decl(self):
        self.eat(TokenType.DO)
        self.eat(TokenType.WHILE)
        do_while = ast.DoWhile()
        do_while.condition = self.expression_decl()
        do_while.block = self.block(ending_block_token_type=TokenType.ENDDO)

        return do_while
    """
       if_decl ::= 'if' condition block ( 'else' block )? 'endif' 
    """
    def if_decl(self):
        self.eat(TokenType.IF)
        if_stmt = ast.IfStatement()
        if_stmt.condition = self.expression_decl()
        if_stmt.consequence = ast.Block()

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
        assignment_decl ::= identifier '=' expression_decl LBREAK
    """
    def assignment_decl(self):
        ident = self.identifier()
        self.eat(TokenType.ASSIGN)
        value = self.expression_decl()
        return ast.Assignment(token=ident, value=value)
    """
        return_decl ::= 'return' (expression_decl)?
    """
    def return_decl(self):
        return_stmt = ast.ReturnStmt()
        self.eat(TokenType.RETURN)

        if self.cur_token.type != TokenType.LBREAK:
            return_stmt.value = self.expression_decl()
        else:
            return_stmt.value = ast.Boolean(value=True)

        return return_stmt
    """
        print_decl ::= 'print' expression_decl (',' expression_decl)*
    """
    def print_decl(self):
        print_stmt = ast.PrintStmt()
        self.eat(TokenType.PRINT)
        print_stmt.arguments.append(self.expression_decl())

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            print_stmt.arguments.append(self.expression_decl())

        return print_stmt

    """
    expression_decl ::= logic_or
    logic_or         ::= logic_and ('or' logic_and )*
    logic_and        ::= equality ('and' equality )*
    equality        ::= comparison ( ('==' | '!=') comparison )*
    comparison      ::= term ( ( '<' | '>' | '<=' | '>=' ) term )*
    term            ::= factor ( ( '+' | '-' ) factor )*
    factor          ::= unary ( ( '*' | '/' ) unary )*
    unary           ::= primary ( ( '!' | '-') primary )*
    primary         ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER | '(' expression_decl ')'  
    """
    def expression_decl(self):
        exp = self.logic_or()

        if self.cur_token.type == TokenType.LBREAK:
            self.eat(TokenType.LBREAK)

        return exp
    """
    logic_or ::= logic_and ('or' logic_and )* 
    """
    def logic_or(self):
        or_exp = self.logic_and()
        while self.cur_token.type == TokenType.OR:
            tok = self.cur_token
            self.eat(TokenType.OR)
            or_exp = ast.BinaryOp(
                left=or_exp,
                operator=tok.value,
                right=self.logic_and(),
            )

        return or_exp
    """
    logic_and ::= equality ('and' equality )*  
    """
    def logic_and(self):
        and_exp = self.equality()
        while self.cur_token.type == TokenType.AND:
            tok = self.cur_token
            self.eat(TokenType.AND)
            and_exp = ast.BinaryOp(
                left=and_exp,
                operator=tok.value,
                right=self.equality(),
            )

        return and_exp
    """
    equality ::= comparison ( ('==' | '!=') comparison )*  
    """
    def equality(self):
        comp_exp = self.comparison()
        while self.cur_token.type in (TokenType.EQUAL, TokenType.NOT_EQUAL):
            tok = self.cur_token
            self.eat(tok.type)  # '==' ó '!='
            comp_exp = ast.BinaryOp(
                left=comp_exp,
                operator=tok.value,
                right=self.comparison(),
            )

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
            term_exp = ast.BinaryOp(
                left=term_exp,
                operator=tok.value,
                right=self.term(),
            )

        return term_exp
    """
    term ::= factor ( ( '+' | '-' ) factor )*   
    """
    def term(self):
        factor_exp = self.factor()
        while self.cur_token.type in (TokenType.PLUS, TokenType.MINUS):
            tok = self.cur_token
            self.eat(tok.type)  # '+' ó '-'
            factor_exp = ast.BinaryOp(
                left=factor_exp,
                operator=tok.value,
                right=self.factor(),
            )

        return factor_exp
    """
    factor ::= unary ( ( '*' | '/' ) unary )*
    """
    def factor(self):
        unary_exp = self.unary()
        while self.cur_token.type in (TokenType.MUL, TokenType.DIV):
            tok = self.cur_token
            self.eat(tok.type)  # '*' ó '/'
            unary_exp = ast.BinaryOp(
                left=unary_exp,
                operator=tok.value,
                right=self.unary(),
            )

        return unary_exp
    """
    unary ::= ( '!' | '-') primary | call  
    """
    def unary(self):
        if self.cur_token.type in (TokenType.NOT, TokenType.MINUS):
            operator = self.cur_token.value
            self.eat(self.cur_token.type)  # '!' ó '-'
            return ast.UnaryOp(operator=operator, right=self.unary())
        else:
            return self.function_call()
    """
    functionCall ::= primary ( '(' arguments? ')' )?   
    """
    def function_call(self):
        primary = self.primary()
        if self.cur_token.type == TokenType.LPAREN:
            function_name = primary
            primary = ast.FunctionCall()
            primary.name = function_name

            self.eat(TokenType.LPAREN)

            if self.cur_token.type != TokenType.RPAREN:
                primary.arguments = self.parse_arguments()

            self.eat(TokenType.RPAREN)

        return primary
    """
    primary ::= '.t.' | '.f.' | '.null.' | NUMBER | STRING | IDENTIFIER | '(' expression_decl ')'  
    """
    def primary(self):
        tok = self.cur_token

        if tok.type in (TokenType.TRUE, TokenType.FALSE):
            self.eat(tok.type)  # .T. ó .F.
            return ast.Boolean(value=(tok.type == TokenType.TRUE))

        elif tok.type == TokenType.NULL:
            self.eat(TokenType.NULL)
            return ast.Null()

        elif tok.type == TokenType.INT:
            self.eat(TokenType.INT)
            return ast.Integer(value=int(tok.value))

        elif tok.type == TokenType.STRING:
            self.eat(TokenType.STRING)
            return ast.String(value=tok.value)

        elif tok.type == TokenType.IDENT:
            return self.identifier()

        elif tok.type == TokenType.LPAREN:
            self.eat(TokenType.LPAREN)
            exp = self.expression_decl()
            self.eat(TokenType.RPAREN)
            return exp
    """
    parse_arguments ::= expression_decl ( ',' expression_decl )*  
    """
    def parse_arguments(self):
        args = [self.expression_decl()]

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            args.append(self.expression_decl())

        return args
    """
    parse_parameters ::= identifier ( ',' identifier )*  
    """
    def parse_parameters(self):
        identifiers = [self.identifier()]

        while self.cur_token.type == TokenType.COMMA:
            self.eat(TokenType.COMMA)
            identifiers.append(self.identifier())

        return identifiers

    def identifier(self):
        name = self.cur_token.value
        self.eat(TokenType.IDENT)
        return ast.Identifier(value=name)
