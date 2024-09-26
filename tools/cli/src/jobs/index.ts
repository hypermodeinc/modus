export class JobBase {
    /**
     * Initialize the job
     * @returns void
     */
    init(): void {}
    /**
     * Returns true if ok, false if error
     * @returns boolean
     */
    run(): boolean {
        return false;
    }
}