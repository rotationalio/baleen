export function numberFormat(nb: number | bigint): string {
    return new Intl.NumberFormat().format(nb);
}
