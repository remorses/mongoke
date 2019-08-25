
# where_filter, collection
@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = ${{repr_eval_dict(where_filter, '    ')}}
    x = await ctx['db']['${{collection}}'].find_one(where)
${{repr_guards_after_checks(guards_after, '    ')}}
${{repr_disambiguations(disambiguations, '    ')}}
    return x