int i = 0;
2a
char c = 'asd'
string s = "fasgasdr
'sad
0xfgh
085
0b12
_fag_gsdg
/*gh
gsda
hsdf
fasd
dfas
**/
//****
typedef struct NODE_
{
    struct NODE_ *next;
    int val;
    int fre;
} node;
node *createHead()
{
    node *h = (node *)malloc(sizeof(node));
    h->next = NULL;
    h->val = -1;
    h->fre = 0;
    return h;
}
node *createNode()
{
    node *temp = (node *)malloc(sizeof(node));
    temp->fre = 0;
    scanf("%d", &temp->val);
    return temp;
}
void printNode(node *head)
{
    ++i;
    printf("%d:\n", i);
    node *temp = head->next;
    head = head->next;
    while (head)
    {
        printf("%d ", head->val);
        head = head->next;
    }
    printf("\n");
    while (temp)
    {
        printf("%d ", temp->fre);
        temp = temp->next;
    }
    printf("\n");
}
int main()
{
    node *head = createHead();
    node *temp = head;
    int a[10] = {0, 1, 2, 4, 4, 4, 6, 6, 9, 9};
    for (int i = 0; i < 10; i++)
    {
        head->next = createNode();
        head = head->next;
        head->next = NULL;
    }
    for (int i = 0; i < 10; i++)
    {
        findNode(temp, a[i]);
    }
}