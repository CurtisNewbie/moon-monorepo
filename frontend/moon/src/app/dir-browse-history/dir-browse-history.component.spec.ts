import { NavType } from "../routes";
import {
  DirBrowseHistoryComponent,
  DirBrowseRecord,
} from "./dir-browse-history.component";

describe("DirBrowseHistoryComponent resume", () => {
  let navigateTo: jasmine.Spy;
  let component: DirBrowseHistoryComponent;

  beforeEach(() => {
    navigateTo = jasmine.createSpy("navigateTo");
    component = new DirBrowseHistoryComponent(
      { open: () => undefined } as any,
      {} as any,
      { navigateTo } as any
    );
  });

  it("resumes under the entry's own sort and file", () => {
    const entry: DirBrowseRecord = {
      dirKey: "dir-1",
      name: "my dir",
      fileKey: "file-9",
      orderBy: "time",
      orderDir: "desc",
      time: 123,
    };

    component.resumeReading(entry);

    expect(navigateTo).toHaveBeenCalledWith(NavType.MANAGE_FILES, [
      {
        parentDirKey: "dir-1",
        targetFileKey: "file-9",
        orderBy: "time",
        orderDir: "desc",
        autoPreview: true,
      },
    ]);
  });

  it("does not force the fixed name sort", () => {
    component.resumeReading({
      dirKey: "dir-2",
      name: "other dir",
      fileKey: "file-3",
      orderBy: "custom",
      orderDir: "asc",
      time: 456,
    });

    const params = navigateTo.calls.mostRecent().args[1][0];
    expect(params.orderBy).toBe("custom");
    expect(params.orderBy).not.toBe("name");
  });
});
