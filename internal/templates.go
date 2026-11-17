package internal

type CsprojTempSrc struct {
	NetVersion string
}

const CsprojTemp = `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>{{.NetVersion}}</TargetFramework>
	<ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
	<NoWarn>CS8604</NoWarn>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="Npgsql" Version="10.0.3" />
    <PackageReference Include="Npgsql.DependencyInjection" Version="10.0.3" />
  </ItemGroup>
</Project>
`

type ExtensionsCsSrc struct {
	Namespace string
	ClassName string
	FuncName  string
}

const ExtensionsCsTemp = `//autogen dont change
namespace {{.Namespace}};
using Microsoft.Extensions.DependencyInjection;

public static class Extensions
{
    public static IServiceCollection {{.FuncName}}(this IServiceCollection services, string connectionString)
    {
        services.AddNpgsqlDataSource(connectionString);
		services.AddScoped<{{.ClassName}}>();

        return services;
    }
}`

type DbCtxSrc struct {
	Namespace string
	Classname string
}

const DbCtxTemp = `//autogen dont change
using Npgsql;

namespace {{.Namespace}};

public interface IDbCtx
{
    public NpgsqlCommand BuildCommand();
}

public class {{.Classname}} : IDbCtx
{
    NpgsqlDataSource _ds;

    public {{.Classname}}(NpgsqlDataSource ds)
    {
        _ds = ds;
    }

    public NpgsqlCommand BuildCommand()
    {
        return _ds.CreateCommand();
    }

    public async Task<Transaction> BeginTx()
    {
        var connection = await _ds.OpenConnectionAsync();
        var tx = await connection.BeginTransactionAsync();

        return new Transaction(connection, tx);
    }
}

public class Transaction : IDbCtx, IAsyncDisposable
{
    NpgsqlConnection _connection;
    NpgsqlTransaction _tx;

    public Transaction(NpgsqlConnection c, NpgsqlTransaction t)
    {
        _connection = c;
        _tx = t;
    }

    public NpgsqlCommand BuildCommand()
    {
        return new NpgsqlCommand(null, _connection, _tx);
    }

    public async Task Commit()
    {
        await _tx.CommitAsync();
        await _connection.CloseAsync();
    }

	public async Task Rollback()
	{
        await _tx.DisposeAsync();
        await _connection.DisposeAsync();
	}

    public async ValueTask DisposeAsync()
    {
        await _tx.DisposeAsync();
        await _connection.DisposeAsync();
    }
}`

type QueryFileSrc struct {
	Namespace string
	Classname string
	Queries   []string
}

const QueryFileTemp = `//autogen dont change
namespace {{.Namespace}};

public partial class Query
{
	IDbCtx _ctx;
    public Query(IDbCtx ctx)
    {
        _ctx = ctx;
    }

    object? ResolveValue(object? value) => value == null ? DBNull.Value : value;

{{range .Queries}} {{.}} {{end}}
}
`

type QuerySrc struct {
	Name   string
	Text   string
	Params []Field
	Result []Field
}

const OneQueryTemp = `
{{- define "nettype"}}{{.NetType}}{{if .IsArray}}[]{{end}}{{if .IsNullable}}?{{end}}{{end -}}
{{- define "field"}}{{template "nettype" .}} {{.NetRecordName}}{{end -}}
{{- define "param"}}{{template "nettype" .}} {{.NetArgName}}{{end -}}
{{- define "field_read"}}{{if .IsNullable}}reader.IsDBNull({{.Idx}}) ? null : reader.GetFieldValue<{{template "nettype" .}}>({{.Idx}}){{else}}reader.GetFieldValue<{{template "nettype" .}}>({{.Idx}}){{end}}{{end -}}
    readonly string {{.Name}}Query = @"
{{.Text}}";
	public record {{.Name}}Record({{range $i,$f := .Result}}{{if $i}},{{end}}{{template "field" $f}}{{end}});
    public async Task<{{.Name}}Record?> {{.Name}}({{range $i,$f := .Params}}{{if $i}},{{end}}{{template "param" $f}}{{end}})
    {
        await using var cmd = _ctx.BuildCommand();
        cmd.CommandText = {{.Name}}Query;

		{{range .Params}}cmd.Parameters.AddWithValue(ResolveValue({{.NetArgName}}));
		{{end}}
        await using var reader = await cmd.ExecuteReaderAsync();

        var ok = await reader.ReadAsync();
        if (!ok) return null;

        return new {{.Name}}Record(
			{{range .Result}}{{if .Idx}},
			{{end}}{{template "field_read" .}}{{end}});
    }

`

const ManyQueryTemp = `
{{- define "nettype"}}{{.NetType}}{{if .IsArray}}[]{{end}}{{if .IsNullable}}?{{end}}{{end -}}
{{- define "field"}}{{template "nettype" .}} {{.NetRecordName}}{{end -}}
{{- define "param"}}{{template "nettype" .}} {{.NetArgName}}{{end -}}
{{- define "field_read"}}{{if .IsNullable}}reader.IsDBNull({{.Idx}}) ? null : reader.GetFieldValue<{{template "nettype" .}}>({{.Idx}}){{else}}reader.GetFieldValue<{{template "nettype" .}}>({{.Idx}}){{end}}{{end -}}
    readonly string {{.Name}}Query = @"
{{.Text}}";
    public record {{.Name}}Record({{range $i,$f := .Result}}{{if $i}},{{end}}{{template "field" $f}}{{end}});
    public async Task<IEnumerable<{{.Name}}Record>> {{.Name}}({{range $i,$f := .Params}}{{if $i}},{{end}}{{template "param" $f}}{{end}})
    {
        await using var cmd = _ctx.BuildCommand();
        cmd.CommandText = {{.Name}}Query;
		{{range .Params}}cmd.Parameters.AddWithValue(ResolveValue({{.NetArgName}}));
		{{end}}
        await using var reader = await cmd.ExecuteReaderAsync();

        var result = new List<{{.Name}}Record>();
        while (await reader.ReadAsync())
        {
			result.Add(new {{.Name}}Record(
				{{range .Result}}{{if .Idx}},
				{{end}}{{template "field_read" .}}{{end}}));
        }

        return result;
    }

`

const ExecQueryTemp = `
{{- define "nettype"}}{{.NetType}}{{if .IsArray}}[]{{end}}{{if .IsNullable}}?{{end}}{{end -}}
{{- define "param"}}{{template "nettype" .}} {{.NetArgName}}{{end -}}
    readonly string {{.Name}}Query = @"
{{.Text}}";
    public async Task<int> {{.Name}}({{range $i,$f := .Params}}{{if $i}},{{end}}{{template "param" $f}}{{end}})
    {
        await using var cmd = _ctx.BuildCommand();
        cmd.CommandText = {{.Name}}Query;
		{{range .Params}}cmd.Parameters.AddWithValue(ResolveValue({{.NetArgName}}));
		{{end}}
        return await cmd.ExecuteNonQueryAsync();
    }

`
